package v2

import (
	"github.com/BurntSushi/toml"
	shim_config "github.com/confidential-containers/enclave-cc/src/shim/config"
)

var (
	RuntimeClass string
)

func parseConfig(path string) error {
	var cfg shim_config.Config

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return err
	}

	RuntimeClass = cfg.RuntimeClass

	return nil
}
