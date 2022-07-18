package v2

import (
	"context"
	"os"
	"path/filepath"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/oci"
	shimtypes "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/types"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
)

func create(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest) (*runc.Container, error) {
	ociSpec, err := config.LoadSpec(filepath.Join(r.Bundle, "config.json"))
	if err != nil {
		return nil, err
	}
	containerType, err := oci.ContainerType(*ociSpec)
	if err != nil {
		return nil, err
	}

	container, err := runc.NewContainer(ctx, s.platform, r)
	if err != nil {
		return nil, err
	}

	switch containerType {
	case shimtypes.PodSandbox:
		sandboxNamspace, err := oci.SandboxNamespace(*ociSpec)
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
