// The codebase is inherited from kata-containers with the modifications.

package v2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/config"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
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
	testNamespaceTypeJubesystem    = "kube-system"
)

const (
	testDirMode  = os.FileMode(0750)
	testFileMode = os.FileMode(0640)

	busyboxConfigJson = `
{
	"ociVersion": "1.0.1-dev",
	"process": {
		"terminal": true,
		"user": {
			"uid": 0,
			"gid": 0
		},
		"args": [
			"sh"
		],
		"env": [
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm"
		],
		"cwd": "/",
		"capabilities": {
			"bounding": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"effective": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"inheritable": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"permitted": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"ambient": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			]
		},
		"rlimits": [
			{
				"type": "RLIMIT_NOFILE",
				"hard": 1024,
				"soft": 1024
			}
		],
		"noNewPrivileges": true
	},
	"root": {
		"path": "rootfs",
		"readonly": true
	},
	"hostname": "runc",
	"mounts": [
		{
			"destination": "/proc",
			"type": "proc",
			"source": "proc"
		},
		{
			"destination": "/dev",
			"type": "tmpfs",
			"source": "tmpfs",
			"options": [
				"nosuid",
				"strictatime",
				"mode=755",
				"size=65536k"
			]
		},
		{
			"destination": "/dev/pts",
			"type": "devpts",
			"source": "devpts",
			"options": [
				"nosuid",
				"noexec",
				"newinstance",
				"ptmxmode=0666",
				"mode=0620",
				"gid=5"
			]
		},
		{
			"destination": "/dev/shm",
			"type": "tmpfs",
			"source": "shm",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"mode=1777",
				"size=65536k"
			]
		},
		{
			"destination": "/dev/mqueue",
			"type": "mqueue",
			"source": "mqueue",
			"options": [
				"nosuid",
				"noexec",
				"nodev"
			]
		},
		{
			"destination": "/sys",
			"type": "sysfs",
			"source": "sysfs",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"ro"
			]
		},
		{
			"destination": "/sys/fs/cgroup",
			"type": "cgroup",
			"source": "cgroup",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"relatime",
				"ro"
			]
		}
	],
	"linux": {
		"resources": {
			"devices": [
				{
					"allow": false,
					"access": "rwm"
				}
			]
		},
		"namespaces": [
			{
				"type": "pid"
			},
			{
				"type": "network"
			},
			{
				"type": "ipc"
			},
			{
				"type": "uts"
			},
			{
				"type": "mount"
			}
		],
		"maskedPaths": [
			"/proc/acpi",
			"/proc/asound",
			"/proc/kcore",
			"/proc/keys",
			"/proc/latency_stats",
			"/proc/timer_list",
			"/proc/timer_stats",
			"/proc/sched_debug",
			"/sys/firmware",
			"/proc/scsi"
		],
		"readonlyPaths": [
			"/proc/bus",
			"/proc/fs",
			"/proc/irq",
			"/proc/sys",
			"/proc/sysrq-trigger"
		]
	}
}`
)

func SetupOCIConfigFile(t *testing.T) (rootPath string, bundlePath, ociConfigFile string) {
	assert := assert.New(t)

	tmpdir := t.TempDir()

	bundlePath = filepath.Join(tmpdir, "bundle")
	err := os.MkdirAll(bundlePath, testDirMode)
	assert.NoError(err)

	ociConfigFile = filepath.Join(bundlePath, "config.json")
	err = os.WriteFile(ociConfigFile, []byte(busyboxConfigJson), testFileMode)
	assert.NoError(err)

	return tmpdir, bundlePath, ociConfigFile
}

func TestCreatePod(t *testing.T) {
	assert := assert.New(t)

	_, bundlePath, ociConfigFile := SetupOCIConfigFile(t)

	spec, err := config.LoadSpec(filepath.Join(bundlePath, "config.json"))
	assert.NoError(err)

	// set expected container type and sandboxID
	spec.Annotations = make(map[string]string)
	spec.Annotations[testContainerTypeAnnotation] = testContainerTypeSandbox
	spec.Annotations[testSandboxIDAnnotation] = testSandboxID
	spec.Annotations[testSandboxNamespaceAnnotation] = testNamespaceTypeJubesystem

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

	expectedErr := fmt.Sprintf("OCI runtime create failed: rootfs (%s/rootfs) does not exist", bundlePath)
	assert.Error(err)
	assert.Contains(err.Error(), expectedErr)
}

func TestCreateApp(t *testing.T) {
	assert := assert.New(t)

	_, bundlePath, ociConfigFile := SetupOCIConfigFile(t)

	spec, err := config.LoadSpec(filepath.Join(bundlePath, "config.json"))
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
	_, err = s.Create(ctx, req)

	expectedErr := "failed to mount rootfs component"
	assert.Error(err)
	assert.Contains(err.Error(), expectedErr)

	assert.DirExists(filepath.Join(ContainerBase, testContainerID, "upper"))
	assert.DirExists(filepath.Join(ContainerBase, testContainerID, "work"))

	err = os.RemoveAll(filepath.Join(ContainerBase, testContainerID))
	assert.NoError(err)
}
