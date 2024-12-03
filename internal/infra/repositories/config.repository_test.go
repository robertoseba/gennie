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
		repo := NewConfigRepository("./config")
		assert.Equal(t, "config/config.json", repo.ConfigFile())
	})

	t.Run("saves config to file", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir("./profiles")
		config.SetConversationCacheTo("./cache.gob")
		config.APIKeys.AnthropicApiKey = "anthropic"
		config.APIKeys.GroqApiKey = "groq"
		config.APIKeys.MaritacaApiKey = "maritaca"
		config.APIKeys.OpenAiApiKey = "openai"
		repo := NewConfigRepository(".")
		assert.True(t, config.IsNew())

		err := repo.Save(config)
		assert.Nil(t, err)
		assert.Equal(t, false, config.IsNew())

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
		repo := NewConfigRepository(".")

		err := repo.Save(config)
		assert.Nil(t, err)

		loadedConfig, err := repo.Load()
		assert.Nil(t, err)
		assert.False(t, loadedConfig.IsNew())

		assert.Equal(t, *config, *loadedConfig)

		os.Remove(repo.ConfigFile())
	})

	t.Run("loads default config if file does not exist and sets as new", func(t *testing.T) {
		repo := NewConfigRepository(".")
		loadedConfig, err := repo.Load()
		assert.Nil(t, err)

		expectedConfig := config.NewConfig()
		expectedConfig.SetConversationCacheTo(".")

		assert.True(t, loadedConfig.IsNew())
		assert.Equal(t, expectedConfig, loadedConfig)
	})

	t.Run("when loading default should set cache dir the same as config dir", func(t *testing.T) {
		repo := NewConfigRepository(".")
		loadedConfig, err := repo.Load()
		assert.Nil(t, err)

		assert.Equal(t, ".", loadedConfig.ConversationCacheDir)
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
