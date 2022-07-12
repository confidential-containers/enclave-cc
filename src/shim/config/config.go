package config

type Config struct {
	RuntimeClass string `toml:"runtime_class"`
	LogLevel     string `toml:"log_level"`
}
