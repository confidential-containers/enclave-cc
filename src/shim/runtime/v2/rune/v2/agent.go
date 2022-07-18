package v2

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	shimconfig "github.com/confidential-containers/enclave-cc/src/shim/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/constants"
	types "github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/pkg/process"
	"github.com/containerd/containerd/runtime/v2/runc"
	taskAPI "github.com/containerd/containerd/runtime/v2/task"
	"github.com/containerd/continuity/fs"
	runcC "github.com/containerd/go-runc"
	"github.com/sirupsen/logrus"
)

const (
	configFilename   = "config.json"
	defaultDirPerm   = 0700
	defaultFilePerms = 0600
	agentIDFile      = "agent-id"
)

// The function creates agent enclave container based on a pre-installed OCI bundle
func createAgentContainer(ctx context.Context, s *service, r *taskAPI.CreateTaskRequest) (*runc.Container, error) {
	dir := filepath.Join(agentContainerRootDir, r.ID)
	upperDir := path.Join(dir, "upper")
	workDir := path.Join(dir, "work")
	destDir := path.Join(dir, "merged")
	for _, dir := range []string{upperDir, workDir, destDir} {
		if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
			return nil, err
		}
	}

	var options []string
	// Set index=off when mount overlayfs
	options = append(options, "index=off")
	options = append(options,
		fmt.Sprintf("lowerdir=%s", filepath.Join(agentContainerPath, "rootfs")),
		fmt.Sprintf("workdir=%s", filepath.Join(workDir)),
		fmt.Sprintf("upperdir=%s", filepath.Join(upperDir)),
	)
	r.Rootfs = append(r.Rootfs, &types.Mount{
		Type:    "overlay",
		Source:  "overlay",
		Options: options,
	})
	r.Bundle = destDir

	fs.CopyFile(filepath.Join(r.Bundle, configFilename), filepath.Join(agentContainerPath, configFilename))

	// Create Stdout and Stderr file for agent enclave container
	r.Stdout = filepath.Join(agentContainerRootDir, r.ID, "stdout")
	r.Stderr = filepath.Join(agentContainerRootDir, r.ID, "stderr")
	for _, file := range []string{r.Stdout, r.Stderr} {
		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	}

	agentContainer, err := runc.NewContainer(ctx, s.platform, r)
	if err != nil {
		return nil, err
	}

	return agentContainer, nil
}

// Cleanup the agent enclave container resource
func cleanupAgentContainer(ctx context.Context, id string) error {
	var cfg shimconfig.Config
	if _, err := toml.DecodeFile(constants.ConfigurationPath, &cfg); err != nil {
		return err
	}
	rootdir := cfg.Containerd.AgentContainerRootDir
	path := filepath.Join(rootdir, id, "merged")

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}

	runtime, err := runc.ReadRuntime(path)
	if err != nil {
		return err
	}

	opts, err := runc.ReadOptions(path)
	if err != nil {
		return err
	}
	root := process.RuncRoot
	if opts != nil && opts.Root != "" {
		root = opts.Root
	}

	logrus.WithFields(logrus.Fields{
		"root":    root,
		"path":    path,
		"ns":      ns,
		"runtime": runtime,
	}).Debug("agent enclave Container Cleanup()")

	r := process.NewRunc(root, path, ns, runtime, "", false)
	if err := r.Delete(ctx, id, &runcC.DeleteOpts{
		Force: true,
	}); err != nil {
		logrus.WithError(err).Warn("failed to remove agent enclave container")
	}
	if err := mount.UnmountAll(filepath.Join(path, "rootfs"), 0); err != nil {
		logrus.WithError(err).Warn("failed to cleanup rootfs mount")
	}
	if err := os.RemoveAll(filepath.Join(rootdir, id)); err != nil {
		logrus.WithError(err).Warn("failed to remove agent enclave container path")
	}

	return nil
}
