// The codebase is inherited from kata-containers with the modifications.

package v2

import (
	"context"
	"fmt"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/image"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/runtime/v2/shim"
	"github.com/containerd/containerd/runtime/v2/task"
	"github.com/containerd/ttrpc"
	"github.com/sirupsen/logrus"
)

func init() {
	plugin.Register(&plugin.Registration{
		Type:     plugin.TTRPCPlugin,
		ID:       "image",
		Requires: []plugin.Type{"*"},
		InitFn:   initImageService,
	})
}

type ImageService struct {
	s *service
}

func initImageService(ic *plugin.InitContext) (interface{}, error) {
	i, err := ic.GetByID(plugin.TTRPCPlugin, "task")
	if err != nil {
		return nil, fmt.Errorf("get task plugin error. %w", err)
	}
	task := i.(shim.TaskService)
	s := task.TaskService.(*service)
	is := &ImageService{s: s}
	return is, nil
}

func (is *ImageService) RegisterTTRPC(server *ttrpc.Server) error {
	task.RegisterImageService(server, is)
	return nil
}

// PullImage and unbundle ready for container creation
func (is *ImageService) PullImage(ctx context.Context, req *task.PullImageRequest) (_ *task.PullImageResponse, err error) {
	is.s.mu.Lock()
	defer is.s.mu.Unlock()

	logrus.WithFields(logrus.Fields{
		"image": req.Image,
	}).Debug("Making image pull request")

	r := &image.PullImageReq{
		Image: req.Image,
	}

	resp, err := is.s.agent.PullImage(ctx, r)
	if err != nil {
		logrus.Errorf("rune runtime PullImage err. %v", err)
		return nil, err
	}

	return &task.PullImageResponse{
		ImageRef: resp.ImageRef,
	}, err
}
