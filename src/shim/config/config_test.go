package config

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestDecodeConfig(t *testing.T) {
	assert := assert.New(t)

	var cfg Config
	text := `log_level = "debug" # "debug" "info" "warn" "error"

[containerd]
    agent_container_instance = "/opt/enclave-cc/agent-instance/"
    agent_container_root_dir = "/run/containerd/agentenclave"
    agent_url = "tcp://0.0.0.0:7788"
    boot_container_instance = "/opt/enclave-cc/boot-instance/"
`
	_, err := toml.Decode(text, &cfg)

	assert.NoError(err)
}
