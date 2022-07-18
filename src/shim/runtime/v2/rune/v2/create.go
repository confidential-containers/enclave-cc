package v2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/oci"
	shimtypes "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/types"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
	"github.com/sirupsen/logrus"
)

const ContainerBase = "/run/enclave-cc/app"

func create(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest) (*runc.Container, error) {
	ociSpec, err := config.LoadSpec(filepath.Join(r.Bundle, "config.json"))
	if err != nil {
		return nil, err
	}
	containerType, err := oci.ContainerType(*ociSpec)
	if err != nil {
		return nil, err
	}
	sandboxNamspace, err := oci.SandboxNamespace(*ociSpec)
	if err != nil {
		return nil, err
	}

	var container *runc.Container

	switch containerType {
	case shimtypes.PodSandbox:
		container, err = runc.NewContainer(ctx, s.platform, r)
		if err != nil {
			return nil, err
		}

		if sandboxNamspace != shimtypes.KubeSystem {
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

			if err := writeAgentIdFile(r.Bundle, ar.ID); err != nil {
				return nil, err
			}

			s.agentID = ar.ID
			s.pauseID = r.ID
			s.containers[ar.ID] = agentContainer
			s.agent = &agent{
				ID:     agentContainer.ID,
				Bundle: agentContainer.Bundle,
				URL:    AgentUrl,
			}
		}
	case shimtypes.PodContainer:
		if sandboxNamspace != shimtypes.KubeSystem {
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
				if err := os.MkdirAll(dir, 0711); err != nil && !os.IsExist(err) {
					return nil, err
				}
			}
			// sefsDir store the unionfs images (based on sefs)
			sefsDir := filepath.Join(agentContainerRootDir, s.agentID, "merged/rootfs/images", cid)

			var options []string
			// Set index=off when mount overlayfs
			options = append(options, "index=off")
			options = append(options,
				fmt.Sprintf("workdir=%s", workDir),
				fmt.Sprintf("upperdir=%s", upperDir),
			)
			options = append(options, fmt.Sprintf("lowerdir=%s:%s", sefsDir, filepath.Join(bootContainerPath, "rootfs")))
			r.Rootfs = append(r.Rootfs, &types.Mount{
				Type:    "overlay",
				Source:  "overlay",
				Options: options,
			})

			// Update the Root.Path field in container spec from
			// "/var/lib/containerd/io.containerd.grpc.v1.cri/containers/<id>/rootfs" to "rootfs"
			// TODO: config.json will be updated by agent enclave container
			err = config.UpdateRootPathConfig(filepath.Join(r.Bundle, "config.json"), "rootfs")
			if err != nil {
				return nil, err
			}

			shimLog.WithFields(logrus.Fields{
				"Rootfs": r.Rootfs,
				"Bundle": r.Bundle,
			}).Debug("Create app enclave container based on sefs image")
		}

		container, err = runc.NewContainer(ctx, s.platform, r)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// ReadAgentIdFile reads the agent container id information from the path
func readAgentIdFile(path string) (string, error) {
	data, err := os.ReadFile(filepath.Join(path, "agent.id"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteAgentIdFile writes the agent container id information into the path
func writeAgentIdFile(path, id string) error {
	return os.WriteFile(filepath.Join(path, "agent.id"), []byte(id), 0600)
}
