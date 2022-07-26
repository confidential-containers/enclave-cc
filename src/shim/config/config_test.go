package config

import (
	_ "embed"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

//go:embed config.toml
var shimConfig string

func TestDecodeConfig(t *testing.T) {
	assert := assert.New(t)

	var cfg Config
	_, err := toml.Decode(shimConfig, &cfg)

	assert.NoError(err)
}
