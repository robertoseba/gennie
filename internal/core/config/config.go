package config

import "time"

type Config struct {
	APIKeys              ProviderApiKeys
	Ollama               OllamaConfig
	IsTerminalPretty     bool
	ConversationCacheDir string
	ProfilesDirPath      string
	HttpTimeout          time.Duration
}

type ProviderApiKeys struct {
	OpenAiApiKey    string
	AnthropicApiKey string
	MaritacaApiKey  string
	GroqApiKey      string
}

type OllamaConfig struct {
	Host  string
	Model string
}

func NewConfig() *Config {
	return &Config{
		APIKeys: ProviderApiKeys{
			OpenAiApiKey:    "",
			AnthropicApiKey: "",
			MaritacaApiKey:  "",
			GroqApiKey:      "",
		},
		Ollama: OllamaConfig{
			Host:  "",
			Model: "",
		},
		IsTerminalPretty:     true,
		ConversationCacheDir: "",
		ProfilesDirPath:      "",
		HttpTimeout:          time.Second * 60,
	}
}

func (c *Config) SetConversationCacheTo(filepath string) {
	c.ConversationCacheDir = filepath
}

func (c *Config) SetProfilesDir(dirPath string) {
	c.ProfilesDirPath = dirPath
}

func (c *Config) SetOllama(host, model string) {
	c.Ollama.Host = host
	c.Ollama.Model = model
}
