package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()

	require.Equal(t, "", c.APIKeys.OpenAiApiKey)
	require.Equal(t, "", c.APIKeys.AnthropicApiKey)
	require.Equal(t, "", c.APIKeys.MaritacaApiKey)
	require.Equal(t, "", c.APIKeys.GroqApiKey)
	require.Equal(t, "", c.Ollama.Host)
	require.Equal(t, "", c.Ollama.Model)
	require.True(t, c.IsTerminalPretty)
	require.Equal(t, "", c.ConversationCacheDir)
	require.Equal(t, "", c.ProfilesDirPath)
	require.InEpsilon(t, float64(60), c.HttpTimeout.Seconds(), 0.1)

	t.Run("Set conversation cache dir", func(t *testing.T) {
		c := NewConfig()
		c.SetConversationCacheTo("cache")
		require.Equal(t, "cache", c.ConversationCacheDir)
	})

	t.Run("SetProfilesDir", func(t *testing.T) {
		c := NewConfig()
		c.SetProfilesDir("profiles")
		require.Equal(t, "profiles", c.ProfilesDirPath)
	})

	t.Run("SetOllama", func(t *testing.T) {
		c := NewConfig()
		c.SetOllama("host", "model")
		require.Equal(t, "host", c.Ollama.Host)
		require.Equal(t, "model", c.Ollama.Model)
	})

	t.Run("Is new", func(t *testing.T) {
		c := NewConfig()
		require.True(t, c.IsNew())

		c.MarkAsNotNew()
		require.False(t, c.IsNew())
	})
}
