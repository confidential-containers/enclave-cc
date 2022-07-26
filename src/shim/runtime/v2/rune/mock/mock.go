// The codebase is inherited from kata-containers with the modifications.

package mock

import (
	"context"
	"fmt"
	"net"
	"net/url"

	pb "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/agent/grpc"

	"github.com/containerd/ttrpc"
)

var testMockAgentSockURLTempl = "tcp://127.0.0.1:8822"

func GenerateMockAgentSock() (string, error) {
	return testMockAgentSockURLTempl, nil
}

// SockTTRPCMock is the ttrpc-based mock sock backend implementation
type SockTTRPCMock struct {
	// SockTTRPCMockImp is the structure implementing
	// the ttrpc interface we want the mock sock server to serve.
	SockTTRPCMockImp

	listener net.Listener
}

func (st *SockTTRPCMock) ttrpcRegister(s *ttrpc.Server) {
	pb.RegisterImageService(s, &st.SockTTRPCMockImp)
}

// Start starts the ttrpc-based mock server
func (st *SockTTRPCMock) Start(socketAddr string) error {
	if socketAddr == "" {
		return fmt.Errorf("missing Socket Address")
	}

	url, err := url.Parse(socketAddr)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", url.Host)
	if err != nil {
		return err
	}

	st.listener = l

	ttrpcServer, err := ttrpc.NewServer()
	if err != nil {
		return err
	}
	st.ttrpcRegister(ttrpcServer)

	go func() {
		ttrpcServer.Serve(context.Background(), l)
	}()

	return nil
}

// Stop stops the ttrpc-based mock server
func (st *SockTTRPCMock) Stop() error {
	if st.listener == nil {
		return fmt.Errorf("missing mock sock listener")
	}

	return st.listener.Close()
}

type SockTTRPCMockImp struct{}

func (p *SockTTRPCMockImp) PullImage(ctx context.Context, req *pb.PullImageRequest) (*pb.PullImageResponse, error) {
	return &pb.PullImageResponse{}, nil
}
