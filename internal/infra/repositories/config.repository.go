package repositories

import (
	"encoding/json"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/core/config"
)

type ConfigRepository struct {
	filename     string
	dirPath      string
	configCached *config.Config
}

func NewConfigRepository(configDir string) *ConfigRepository {
	return &ConfigRepository{
		filename:     "config.json",
		dirPath:      configDir,
		configCached: nil,
	}
}

// Loads the config from a gob file into the Config struct
// If the file does not exist, it returns a new Config with default values
// Once loaded config is cached until the cli quits
func (cr *ConfigRepository) Load() (*config.Config, error) {
	if cr.configCached != nil {
		return cr.configCached, nil
	}

	file := cr.ConfigFile()
	if _, err := os.Stat(file); os.IsNotExist(err) {
		config := config.NewConfig()
		//Default cache dir is the same as config path
		config.ConversationCacheDir = cr.dirPath
		cr.configCached = config
		return config, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config config.Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	cr.configCached = &config
	return &config, nil
}

// Returns the full path to the config file
func (cr *ConfigRepository) ConfigFile() string {
	return path.Join(cr.dirPath, cr.filename)
}

// Saves the config to a json file
// Caches the config in repository cache
func (cr *ConfigRepository) Save(config *config.Config) error {
	config.MarkAsNotNew()
	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(cr.ConfigFile(), content, 0644)
	if err != nil {
		return err
	}
	cr.configCached = config
	return nil
}

// Uses two strategies to define the config directory:
// 1. If the system has a user config directory, it uses it and creates a gennie directory inside it
// 2. If the system does not have a user config directory, it uses the executable directory
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
