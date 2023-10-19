package v2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/oci"
	shimtypes "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/types"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

var ContainerBase = "/run/enclave-cc/app"

func create(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest) (*runc.Container, error) {
	ociSpec, err := config.LoadSpec(filepath.Join(r.Bundle, configFilename))
	if err != nil {
		return nil, err
	}

	containerType, err := oci.ContainerType(*ociSpec)
	if err != nil {
		return nil, err
	}
	sandboxNamespace, err := oci.SandboxNamespace(*ociSpec)
	if err != nil {
		return nil, err
	}

	var container *runc.Container

	switch containerType {
	case shimtypes.PodSandbox:
		container, err = handlePodSandbox(ctx, s, r, sandboxNamespace)
		if err != nil {
			return nil, err
		}
	case shimtypes.PodContainer:
		container, err = handlePodContainer(ctx, s, r, sandboxNamespace, ociSpec)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

func handlePodSandbox(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest, sandboxNamespace string) (*runc.Container, error) {
	container, err := runc.NewContainer(ctx, s.platform, r)
	if err != nil {
		return nil, err
	}

	if sandboxNamespace != shimtypes.KubeSystemNS {
		ar := &taskAPI.CreateTaskRequest{
			ID:       generateID(),
			Terminal: false,
			Options:  r.Options,
		}

		// Create agent enclave container
		agentContainer, err := createAgentContainer(ctx, s, ar)
		if err != nil {
			return nil, err
		}

		if err := writeAgentIDFile(r.Bundle, ar.ID); err != nil {
			return nil, err
		}

		s.agentID = ar.ID
		s.pauseID = r.ID
		s.containers[ar.ID] = agentContainer
		s.agent = &agent{
			ID:     agentContainer.ID,
			Bundle: agentContainer.Bundle,
			URL:    AgentURL,
		}

	}

	return container, nil
}

func handlePodContainer(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest, sandboxNamespace string, ociSpec *specs.Spec) (*runc.Container, error) {
	if sandboxNamespace != shimtypes.KubeSystemNS {
		image, err := oci.GetImage(*ociSpec)
		if err != nil {
			return nil, err
		}
		cid, err := getContainerID(image)
		if err != nil {
			return nil, err
		}

		// Create upperDir and workDir for app container
		upperDir := filepath.Join(ContainerBase, r.ID, "upper")
		workDir := filepath.Join(ContainerBase, r.ID, "work")
		for _, dir := range []string{upperDir, workDir} {
			if err := os.MkdirAll(dir, defaultDirPerm); err != nil && !os.IsExist(err) {
				return nil, err
			}
		}
		// sefsDir store the unionfs images (based on sefs)
		lowerdirs := []string{
			filepath.Join(agentContainerRootDir, s.agentID, "merged/rootfs/images", cid),
			filepath.Join(bootContainerPath, "rootfs"),
		}
		sealDataDir := filepath.Join(agentContainerRootDir, s.agentID, "merged/rootfs/keys", cid)
		if _, err := os.Stat(sealDataDir); !os.IsNotExist(err) {
			lowerdirs = append(lowerdirs, sealDataDir)
		}

		var options []string
		// Set index=off when mount overlayfs
		options = append(options, "index=off")
		options = append(options,
			fmt.Sprintf("workdir=%s", filepath.Join(workDir)),
			fmt.Sprintf("upperdir=%s", filepath.Join(upperDir)),
		)
		options = append(options, fmt.Sprintf("lowerdir=%s", strings.Join(lowerdirs, ":")))
		r.Rootfs = append(r.Rootfs, &types.Mount{
			Type:    "overlay",
			Source:  "overlay",
			Options: options,
		})

		// Update the Root.Path field in container spec from
		// "/var/lib/containerd/io.containerd.grpc.v1.cri/containers/<id>/rootfs" to "rootfs"
		// TODO: config.json will be updated by agent enclave container
		err = config.UpdateRootPathConfig(filepath.Join(r.Bundle, configFilename), "rootfs")
		if err != nil {
			return nil, err
		}

		logrus.WithFields(logrus.Fields{
			"Rootfs": r.Rootfs,
			"Bundle": r.Bundle,
		}).Debug("Create app enclave container based on sefs image")
	}

	container, err := runc.NewContainer(ctx, s.platform, r)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// readAgentIDFile reads the agent container id information from the path
func readAgentIDFile(path string) (string, error) {
	data, err := os.ReadFile(filepath.Join(path, agentIDFile))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// writeAgentIDFile writes the agent container id information into the path
func writeAgentIDFile(path, id string) error {
	return os.WriteFile(filepath.Join(path, agentIDFile), []byte(id), defaultFilePerms)
}
