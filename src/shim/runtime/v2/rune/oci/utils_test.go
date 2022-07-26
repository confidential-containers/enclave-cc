package oci

import (
	"testing"

	shimtypes "github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/types"
	ctrAnnotations "github.com/containerd/containerd/pkg/cri/annotations"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
)

func TestSandboxNamespace(t *testing.T) {
	assert := assert.New(t)

	spec := specs.Spec{}
	_, err := SandboxNamespace(spec)
	assert.Error(err)

	spec = specs.Spec{
		Annotations: map[string]string{
			ctrAnnotations.SandboxNamespace: shimtypes.KubeSystemNS,
		},
	}
	sandboxNamespace, err := SandboxNamespace(spec)
	assert.Equal(shimtypes.KubeSystemNS, sandboxNamespace)
	assert.NoError(err)
}

func TestGetImage(t *testing.T) {
	assert := assert.New(t)

	spec := specs.Spec{}
	_, err := GetImage(spec)
	assert.Error(err)

	ImageName := "busybox"
	spec = specs.Spec{
		Annotations: map[string]string{
			ctrAnnotations.ImageName: ImageName,
		},
	}
	image, err := GetImage(spec)
	assert.Equal(ImageName, image)
	assert.NoError(err)
}
