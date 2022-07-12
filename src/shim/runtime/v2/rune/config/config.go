package config

import (
	"encoding/json"
	"os"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// LoadSpec loads the specification from the provided path.
func LoadSpec(cPath string) (spec *specs.Spec, err error) {
	cf, err := os.Open(cPath)
	if err != nil {
		return nil, err
	}
	defer cf.Close()
	if err := json.NewDecoder(cf).Decode(&spec); err != nil {
		return nil, err
	}
	return spec, nil
}
