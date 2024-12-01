package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()

	assert.Equal(t, "", c.APIKeys.OpenAiApiKey)
	assert.Equal(t, "", c.APIKeys.AnthropicApiKey)
	assert.Equal(t, "", c.APIKeys.MaritacaApiKey)
	assert.Equal(t, "", c.APIKeys.GroqApiKey)
	assert.Equal(t, "", c.Ollama.Host)
	assert.Equal(t, "", c.Ollama.Model)
	assert.Equal(t, true, c.IsTerminalPretty)
	assert.Equal(t, "", c.CacheFilename)
	assert.Equal(t, "", c.ProfilesDirPath)
	assert.Equal(t, float64(60), c.HttpTimeout.Seconds())

	t.Run("SetCacheDir", func(t *testing.T) {
		c := NewConfig()
		c.SetCacheTo("cache")
		assert.Equal(t, "cache", c.CacheFilename)
	})

	t.Run("SetProfilesDir", func(t *testing.T) {
		c := NewConfig()
		c.SetProfilesDir("profiles")
		assert.Equal(t, "profiles", c.ProfilesDirPath)
	})

	t.Run("SetOllama", func(t *testing.T) {
		c := NewConfig()
		c.SetOllama("host", "model")
		assert.Equal(t, "host", c.Ollama.Host)
		assert.Equal(t, "model", c.Ollama.Model)
	})
}
