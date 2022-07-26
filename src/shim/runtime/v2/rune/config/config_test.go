package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConfigJson = `
{
	"ociVersion": "1.0.1-dev",
	"process": {},
	"root": {
		"path": "rootfs"
	},
	"hostname": "runc",
	"mounts": [],
	"linux": {}
}`

func TestUpdateRootPathConfig(t *testing.T) {
	assert := assert.New(t)

	path := filepath.Join("/tmp", "config.json")
	defer os.Remove(path)

	err := os.WriteFile(path, []byte(testConfigJson), 0644)
	assert.NoError(err)

	rpath := "test-path"
	err = UpdateRootPathConfig(path, rpath)
	assert.NoError(err)

	spec, err := LoadSpec(path)
	assert.NoError(err)
	assert.Equal(spec.Root.Path, rpath)
}
