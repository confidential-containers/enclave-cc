// The codebase is inherited from kata-containers with the modifications.

package v2

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/config"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

const (
	testSandboxID   = "777-77-77777777"
	testContainerID = "42"
	testagentID     = "43"
	testImage       = "docker.io/library/busybox:latest"

	testContainerTypeAnnotation    = "io.kubernetes.cri.container-type"
	testSandboxIDAnnotation        = "io.kubernetes.cri.sandbox-id"
	testSandboxNamespaceAnnotation = "io.kubernetes.cri.sandbox-namespace"
	testImageAnnotation            = "io.kubernetes.cri.image-name"
	testContainerTypeSandbox       = "sandbox"
	testContainerTypeContainer     = "container"
	testNamespaceTypeDefault       = "default"
	testNamespaceTypeKubesystem    = "kube-system"
)

const (
	testDirMode  = os.FileMode(0750)
	testFileMode = os.FileMode(0640)
)

//go:embed testdata/busybox.json
var busyboxConfigJSON []byte

func SetupOCIConfigFile(t *testing.T) (rootPath string, bundlePath, ociConfigFile string) {
	assert := assert.New(t)

	tmpdir := t.TempDir()

	bundlePath = filepath.Join(tmpdir, "bundle")
	err := os.MkdirAll(bundlePath, testDirMode)
	assert.NoError(err)

	ociConfigFile = filepath.Join(bundlePath, configFilename)
	err = os.WriteFile(ociConfigFile, busyboxConfigJSON, testFileMode)
	assert.NoError(err)

	return tmpdir, bundlePath, ociConfigFile
}

func TestCreatePod(t *testing.T) {
	assert := assert.New(t)

	_, bundlePath, ociConfigFile := SetupOCIConfigFile(t)

	spec, err := config.LoadSpec(filepath.Join(bundlePath, configFilename))
	assert.NoError(err)

	// set expected container type and sandboxID
	spec.Annotations = make(map[string]string)
	spec.Annotations[testContainerTypeAnnotation] = testContainerTypeSandbox
	spec.Annotations[testSandboxIDAnnotation] = testSandboxID
	spec.Annotations[testSandboxNamespaceAnnotation] = testNamespaceTypeKubesystem

	// rewrite file
	err = config.SaveSpec(ociConfigFile, spec)
	assert.NoError(err)

	s := &service{
		id:         testContainerID,
		containers: make(map[string]*runc.Container),
		context:    context.Background(),
	}

	req := &taskAPI.CreateTaskRequest{
		ID:       testContainerID,
		Bundle:   bundlePath,
		Terminal: true,
	}

	ctx := namespaces.WithNamespace(context.Background(), "UnitTest")
	_, err = s.Create(ctx, req)

	expectedErr := "OCI runtime create failed"
	assert.Error(err)
	assert.Contains(err.Error(), expectedErr)
}

func TestCreateApp(t *testing.T) {
	assert := assert.New(t)

	_, bundlePath, ociConfigFile := SetupOCIConfigFile(t)

	spec, err := config.LoadSpec(filepath.Join(bundlePath, configFilename))
	assert.NoError(err)

	// set expected container type, sandboxID, SandboxNamespace and Image
	spec.Annotations = make(map[string]string)
	spec.Annotations[testContainerTypeAnnotation] = testContainerTypeContainer
	spec.Annotations[testSandboxIDAnnotation] = testSandboxID
	spec.Annotations[testSandboxNamespaceAnnotation] = testNamespaceTypeDefault
	spec.Annotations[testImageAnnotation] = testImage

	// rewrite file
	err = config.SaveSpec(ociConfigFile, spec)
	assert.NoError(err)

	s := &service{
		id:         testContainerID,
		containers: make(map[string]*runc.Container),
		context:    context.Background(),
		agentID:    testagentID,
	}

	req := &taskAPI.CreateTaskRequest{
		ID:       testContainerID,
		Bundle:   bundlePath,
		Terminal: true,
	}

	ctx := namespaces.WithNamespace(context.Background(), "UnitTest")
	stubs := gostub.Stub(&ContainerBase, t.TempDir())
	defer stubs.Reset()
	_, err = s.Create(ctx, req)

	expectedErr := "failed to mount rootfs component"
	assert.Error(err)
	assert.Contains(err.Error(), expectedErr)
	if t.Failed() {
		t.FailNow()
	}

	assert.DirExists(filepath.Join(ContainerBase, testContainerID, "upper"))
	assert.DirExists(filepath.Join(ContainerBase, testContainerID, "work"))
}
