package config

type ConfigRepository interface {
	Load() (*Config, error)
	Save(config *Config) error
	ConfigFile() string
}
