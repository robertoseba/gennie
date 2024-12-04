package repositories

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/stretchr/testify/require"
)

func TestNewConfigRepository(t *testing.T) {
	t.Run("creates new config repo", func(t *testing.T) {
		repo := NewConfigRepository("./config")
		require.Equal(t, "config/config.json", repo.configFile())
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
		require.True(t, config.IsNew())

		err := repo.Save(config)
		require.NoError(t, err)
		require.False(t, config.IsNew())

		require.FileExists(t, repo.configFile())
		require.NoError(t, err)

		os.Remove(repo.configFile())
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
		require.NoError(t, err)

		loadedConfig, err := repo.Load()
		require.NoError(t, err)
		require.False(t, loadedConfig.IsNew())

		require.Equal(t, *config, *loadedConfig)

		os.Remove(repo.configFile())
	})

	t.Run("loads default config if file does not exist and sets as new", func(t *testing.T) {
		repo := NewConfigRepository(".")
		loadedConfig, err := repo.Load()
		require.NoError(t, err)

		expectedConfig := config.NewConfig()
		expectedConfig.SetConversationCacheTo(".")

		require.True(t, loadedConfig.IsNew())
		require.Equal(t, expectedConfig, loadedConfig)
	})

	t.Run("when loading default should set cache dir the same as config dir", func(t *testing.T) {
		repo := NewConfigRepository(".")
		loadedConfig, err := repo.Load()
		require.NoError(t, err)

		require.Equal(t, ".", loadedConfig.ConversationCacheDir)
	})
}

func TestCreateConfigDir(t *testing.T) {
	t.Run("creates config dir based on XDG_CONFIG_HOME", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", ".") //sets config home to current directory

		dir, err := CreateConfigDir()

		require.NoError(t, err)
		require.Equal(t, "gennie", dir)
		os.Remove(dir)
	})

	t.Run("if already exists, returns the existing dir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", ".") //sets config home to current directory

		os.Mkdir("gennie", 0755)
		dir, err := CreateConfigDir()

		require.NoError(t, err)
		require.Equal(t, "gennie", dir)
		os.Remove(dir)
	})

}
