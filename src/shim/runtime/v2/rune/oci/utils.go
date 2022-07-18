package oci

import (
	"fmt"

	ctrAnnotations "github.com/containerd/containerd/pkg/cri/annotations"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// ContainerType returns the type of container and if the container type was
// found from CRI server's annotations in the container spec.
func ContainerType(spec specs.Spec) (string, error) {
	containerType, ok := spec.Annotations[ctrAnnotations.ContainerType]
	if !ok {
		return "", fmt.Errorf("unknown container type in annotation")
	}

	return containerType, nil
}

// SandboxNamespaceType returns the namespace of sandbox and if the namespace was
// found from CRI server's annotations in the container spec.
func SandboxNamespace(spec specs.Spec) (string, error) {
	sandboxNamespaceType, ok := spec.Annotations[ctrAnnotations.SandboxNamespace]
	if !ok {
		return "", fmt.Errorf("unknown sandbox namespace in annotation")
	}

	return sandboxNamespaceType, nil
}

func GetImage(spec specs.Spec) (string, error) {
	image, ok := spec.Annotations[ctrAnnotations.ImageName]
	if !ok {
		return "", fmt.Errorf("unknown image name in annotation")
	}

	return image, nil
}
