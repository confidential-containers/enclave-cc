package v2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	shim_config "github.com/confidential-containers/enclave-cc/src/shim/config"
)

var (
	RuntimeClass          string
	logLevel              string
	agentContainerRootDir string
	agentContainerPath    string
	AgentURL              string
	bootContainerPath     string
)

func parseConfig(path string) error {
	var cfg shim_config.Config

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return err
	}

	RuntimeClass = cfg.RuntimeClass
	logLevel = cfg.LogLevel
	agentContainerPath = cfg.Containerd.AgentContainerInstance
	agentContainerRootDir = cfg.Containerd.AgentContainerRootDir
	AgentURL = cfg.Containerd.AgentURL
	bootContainerPath = cfg.Containerd.BootContainerInstance

	for _, dir := range []string{agentContainerPath, bootContainerPath} {
		fi, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("pre-installed OCI bundle should exist in %q", dir)
		}
		if !fi.IsDir() {
			return fmt.Errorf("not a directory: %q", dir)
		}
	}

	return nil
}

func generateID() string {
	b := make([]byte, 32)
	rand.Read(b)

	return hex.EncodeToString(b)
}

func getContainerID(image string) (string, error) {
	v := strings.Split(image, "/")
	imageName := v[len(v)-1]
	if len(imageName) <= 0 {
		return "", fmt.Errorf("invalid image name %s", image)
	}

	// ':' have special meaning for umoci during upack in agent enclave container
	cid := strings.Replace(imageName, ":", "_", -1)

	return cid, nil
}
