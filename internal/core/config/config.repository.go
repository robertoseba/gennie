package config

type IConfigRepository interface {
	Load() (*Config, error)
	Save(config Config) error
}
