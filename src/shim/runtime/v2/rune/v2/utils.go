package v2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"

	shim_config "github.com/confidential-containers/enclave-cc/src/shim/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/constants"
)

var (
	logLevel              string
	agentContainerRootDir string
	agentContainerPath    string
)

func parseConfig() error {
	var cfg shim_config.Config

	_, err := toml.DecodeFile(constants.ConfigurationPath, &cfg)
	if err != nil {
		return err
	}

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
