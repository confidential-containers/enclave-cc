package v2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	shim_config "github.com/confidential-containers/enclave-cc/src/shim/config"
)

var (
	RuntimeClass          string
	logLevel              string
	agentContainerRootDir string
	agentContainerPath    string
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

	fi, err := os.Stat(agentContainerPath)
	if err != nil {
		return fmt.Errorf("pre-installed OCI bundle should provided in %s", agentContainerPath)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %s", agentContainerPath)
	}

	return nil
}

func generateID() string {
	b := make([]byte, 32)
	rand.Read(b)

	return hex.EncodeToString(b)
}
