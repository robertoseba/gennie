package repositories

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigRepository(t *testing.T) {
	os.Setenv("XDG_CONFIG_HOME", ".") //sets config home to current directory

	t.Run("creates new config repo", func(t *testing.T) {
		repo, err := NewConfigRepository("./config")

		assert.Nil(t, err)
		assert.Equal(t, "config/gennie_config.gob", repo.ConfigFile())
	})

	t.Run("saves config to file", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir("./profiles")
		config.SetConversationCacheTo("./cache.gob")
		config.APIKeys.AnthropicApiKey = "anthropic"
		config.APIKeys.GroqApiKey = "groq"
		config.APIKeys.MaritacaApiKey = "maritaca"
		config.APIKeys.OpenAiApiKey = "openai"
		repo, err := NewConfigRepository(".")
		assert.Nil(t, err)

		err = repo.Save(config)
		assert.Nil(t, err)

		assert.FileExists(t, repo.ConfigFile())
		assert.Nil(t, err)

		os.Remove(repo.ConfigFile())
	})

	t.Run("loads config from file", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir("./profiles")
		config.SetConversationCacheTo("./cache.gob")
		config.APIKeys.AnthropicApiKey = "anthropic"
		config.APIKeys.GroqApiKey = "groq"
		config.APIKeys.MaritacaApiKey = "maritaca"
		config.APIKeys.OpenAiApiKey = "openai"
		repo, err := NewConfigRepository(".")
		assert.Nil(t, err)

		err = repo.Save(config)
		assert.Nil(t, err)

		loadedConfig, err := repo.Load()
		assert.Nil(t, err)

		assert.Equal(t, *config, *loadedConfig)

		os.Remove(repo.ConfigFile())
	})

	t.Run("loads default config if file does not exist", func(t *testing.T) {
		repo, err := NewConfigRepository(".")
		assert.Nil(t, err)

		loadedConfig, err := repo.Load()
		assert.Nil(t, err)

		assert.Equal(t, *config.NewConfig(), *loadedConfig)
	})
}

func TestCreateConfigDir(t *testing.T) {
	t.Run("creates config dir based on XDG_CONFIG_HOME", func(t *testing.T) {
		os.Setenv("XDG_CONFIG_HOME", ".") //sets config home to current directory
		dir, err := CreateConfigDir()

		assert.Nil(t, err)
		assert.Equal(t, "gennie", dir)
		os.Remove(dir)
	})

	t.Run("if already exists, returns the existing dir", func(t *testing.T) {
		os.Setenv("XDG_CONFIG_HOME", ".") //sets config home to current directory
		os.Mkdir("gennie", 0755)
		dir, err := CreateConfigDir()

		assert.Nil(t, err)
		assert.Equal(t, "gennie", dir)
		os.Remove(dir)
	})
}
