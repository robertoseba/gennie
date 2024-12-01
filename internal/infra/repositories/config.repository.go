package repositories

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/core/config"
)

type ConfigRepository struct {
	filename string
	dirPath  string
}

func NewConfigRepository(configDir string) (*ConfigRepository, error) {
	return &ConfigRepository{
		filename: "gennie_config.gob",
		dirPath:  configDir,
	}, nil
}

func (cr *ConfigRepository) Load() (*config.Config, error) {
	file := cr.ConfigFile()
	if _, err := os.Stat(file); os.IsNotExist(err) {
		config := config.NewConfig()
		return config, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config config.Config
	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (cr *ConfigRepository) ConfigFile() string {
	return path.Join(cr.dirPath, cr.filename)
}

func (cr *ConfigRepository) Save(config *config.Config) error {
	file, err := os.Create(cr.ConfigFile())
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(*config); err != nil {
		return err
	}
	return nil
}

func CreateConfigDir() (string, error) {
	systemConfigDir, err := os.UserConfigDir()

	if err != nil || systemConfigDir == "" {
		return fallbackExecDir()
	}

	defaultConfigDir := path.Join(systemConfigDir, "gennie")

	if _, err := os.Stat(defaultConfigDir); os.IsNotExist(err) {
		err = os.Mkdir(defaultConfigDir, 0755)

		if err != nil {
			return fallbackExecDir()
		}
	}

	return defaultConfigDir, nil
}

func fallbackExecDir() (string, error) {
	filepath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return path.Join(path.Dir(filepath)), nil
}
