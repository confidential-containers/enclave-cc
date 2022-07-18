package v2

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"

	shim_config "github.com/confidential-containers/enclave-cc/src/shim/config"
	"github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/constants"
)

var (
	logLevel              string
	agentContainerRootDir string
	agentContainerPath    string
	AgentUrl              string
	bootContainerPath     string
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
	AgentUrl = cfg.Containerd.AgentUrl
	bootContainerPath = cfg.Containerd.BootContainerInstance

	for _, dir := range []string{agentContainerPath, bootContainerPath} {
		fi, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("pre-installed OCI bundle should provided in %s", dir)
		}
		if !fi.IsDir() {
			return fmt.Errorf("not a directory: %s", dir)
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
	if len(v[len(v)-1]) <= 0 {
		return "", fmt.Errorf("invalid image name %s", image)
	}

	// ':' have special meaning for umoci during upack in agent enclave container
	cid := strings.Replace(v[len(v)-1], ":", "_", -1)

	return cid, nil
}
