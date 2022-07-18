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

func SaveSpec(cPath string, spec *specs.Spec) error {
	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(cPath, data, 0644)
}

func UpdateRootPathConfig(cPath string, rPath string) error {
	spec, err := LoadSpec(cPath)
	if err != nil {
		return err
	}

	spec.Root.Path = rPath

	if err := SaveSpec(cPath, spec); err != nil {
		return err
	}

	return nil
}
