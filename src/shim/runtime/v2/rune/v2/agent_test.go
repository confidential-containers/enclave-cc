// The codebase is inherited from kata-containers with the modifications.

package v2

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/image"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/mock"
	"github.com/stretchr/testify/assert"
)

func TestAgentConnect(t *testing.T) {
	assert := assert.New(t)

	url, err := mock.GenerateMockAgentSock()
	assert.NoError(err)

	SockTTRPCMock := mock.SockTTRPCMock{}
	err = SockTTRPCMock.Start(url)
	assert.NoError(err)
	defer SockTTRPCMock.Stop()

	k := &agent{
		URL: url,
	}

	err = k.connect(context.Background())
	assert.NoError(err)
	assert.NotNil(k.client)
}

func TestAgentDisconnect(t *testing.T) {
	assert := assert.New(t)

	url, err := mock.GenerateMockAgentSock()
	assert.NoError(err)

	SockTTRPCMock := mock.SockTTRPCMock{}
	err = SockTTRPCMock.Start(url)
	assert.NoError(err)
	defer SockTTRPCMock.Stop()

	k := &agent{
		URL: url,
	}

	assert.NoError(k.connect(context.Background()))
	assert.NoError(k.disconnect(context.Background()))
	assert.Nil(k.client)
}

func TestAgentSendReq(t *testing.T) {
	assert := assert.New(t)

	url, err := mock.GenerateMockAgentSock()
	assert.NoError(err)

	SockTTRPCMock := mock.SockTTRPCMock{}
	err = SockTTRPCMock.Start(url)
	assert.NoError(err)
	defer SockTTRPCMock.Stop()

	bundle, err := os.MkdirTemp("", "bundle-test")
	assert.NoError(err)
	defer os.RemoveAll(bundle)

	k := &agent{
		URL:    url,
		Bundle: bundle,
	}

	req := &image.PullImageReq{
		Image: "busybox",
	}

	_, err = k.PullImage(context.Background(), req)
	assert.NoError(err)
	assert.DirExists(filepath.Join(k.Bundle, "rootfs/images/busybox/sefs/lower"))
	assert.DirExists(filepath.Join(k.Bundle, "rootfs/images/busybox/sefs/upper"))
}
