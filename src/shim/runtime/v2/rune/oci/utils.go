package oci

import (
	"errors"
	"fmt"

	ctrAnnotations "github.com/containerd/containerd/pkg/cri/annotations"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var ErrOCIAnnotation = errors.New("OCI Annotation error")

// ContainerType returns the type of container and if the container type was
// found from CRI server's annotations in the container spec.
func ContainerType(spec specs.Spec) (string, error) {
	containerType, ok := spec.Annotations[ctrAnnotations.ContainerType]
	if !ok {
		return "", fmt.Errorf("%w: unknown container type in annotation", ErrOCIAnnotation)
	}

	return containerType, nil
}

// SandboxNamespaceType returns the namespace of sandbox and if the namespace was
// found from CRI server's annotations in the container spec.
func SandboxNamespace(spec specs.Spec) (string, error) {
	sandboxNamespaceType, ok := spec.Annotations[ctrAnnotations.SandboxNamespace]
	if !ok {
		return "", fmt.Errorf("%w: unknown sandbox namespace in annotation", ErrOCIAnnotation)
	}

	return sandboxNamespaceType, nil
}

func GetImage(spec specs.Spec) (string, error) {
	image, ok := spec.Annotations[ctrAnnotations.ImageName]
	if !ok {
		return "", fmt.Errorf("%w: unknown image name in annotation", ErrOCIAnnotation)
	}

	return image, nil
}
