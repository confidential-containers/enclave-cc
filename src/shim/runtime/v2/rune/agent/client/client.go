// The codebase is inherited from kata-containers with the modifications.

package client

import (
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	agentgrpc "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/agent/grpc"
	"github.com/containerd/ttrpc"
)

var defaultDialTimeout = 30 * time.Second

var agentClientFields = logrus.Fields{
	"name":   "agent-client",
	"pid":    os.Getpid(),
	"source": "agent-client",
}

const (
	TCPSocketScheme = "tcp"
)

var agentClientLog = logrus.WithFields(agentClientFields)

// AgentClient is an agent gRPC client connection wrapper for agentgrpc.AgentServiceClient
type AgentClient struct {
	ImageServiceClient agentgrpc.ImageService
	conn               *ttrpc.Client
}

// NewAgentClient creates a new agent gRPC client and handles tcp address.
//
// Supported sock address format is:
//   - tcp://<ip>:<port>
func NewAgentClient(sock string, timeout uint32) (*AgentClient, error) {
	addr, err := url.Parse(sock)
	if err != nil {
		return nil, err
	}
	if addr.Scheme != TCPSocketScheme {
		return nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid scheme: %s, only support tcp scheme", sock)
	}
	if addr.Host == "" {
		return nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid tcp sock scheme: %s", sock)
	}
	grpcAddr := TCPSocketScheme + ":" + addr.Host

	dialTimeout := defaultDialTimeout
	if timeout > 0 {
		dialTimeout = time.Duration(timeout) * time.Second
		agentClientLog.WithField("timeout", timeout).Debug("custom dialing timeout has been set")
	}

	var conn net.Conn
	conn, err = TCPDialer(grpcAddr, dialTimeout)
	if err != nil {
		return nil, err
	}

	client := ttrpc.NewClient(conn)

	return &AgentClient{
		ImageServiceClient: agentgrpc.NewImageClient(client),
		conn:               client,
	}, nil
}

// Close an existing connection to the agent gRPC server.
func (c *AgentClient) Close() error {
	return c.conn.Close()
}

// This would bypass the grpc dialer backoff strategy and handle dial timeout
// internally. Because we do not have a large number of concurrent dialers,
// it is not reasonable to have such aggressive backoffs which would kill
// containers boot up speed. For more information, see
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md
func commonDialer(timeout time.Duration, dialFunc func() (net.Conn, error), timeoutErrMsg error) (net.Conn, error) {
	t := time.NewTimer(timeout)
	cancel := make(chan bool)
	ch := make(chan net.Conn)
	go func() {
		for {
			select {
			case <-cancel:
				// canceled or channel closed
				return
			default:
			}

			conn, err := dialFunc()
			if err == nil {
				// Send conn back iff timer is not fired
				// Otherwise there might be no one left reading it
				if t.Stop() {
					ch <- conn
				} else {
					conn.Close()
				}
				return
			}
		}
	}()

	var conn net.Conn
	var ok bool
	select {
	case conn, ok = <-ch:
		if !ok {
			return nil, timeoutErrMsg
		}
	case <-t.C:
		cancel <- true
		return nil, timeoutErrMsg
	}

	return conn, nil
}

func TCPDialer(sock string, timeout time.Duration) (net.Conn, error) {
	sock = strings.TrimPrefix(sock, "tcp:")

	dialFunc := func() (net.Conn, error) {
		return net.DialTimeout("tcp", sock, timeout)
	}

	timeoutErr := grpcStatus.Errorf(codes.DeadlineExceeded, "timed out connecting to tcp socket %s", sock)
	return commonDialer(timeout, dialFunc, timeoutErr)
}
