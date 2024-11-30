package config

type Config struct {
	APIKeys          ProviderApiKeys
	Ollama           OllamaConfig
	IsTerminalPretty bool
	CacheDirPath     string
	ProfilesDirPath  string
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

func NewConfig() Config {
	return Config{
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
		IsTerminalPretty: true,
		CacheDirPath:     "",
		ProfilesDirPath:  "",
	}
}

func (c *Config) SetCacheDir(path string) {
	c.CacheDirPath = path
}

func (c *Config) SetProfilesDir(path string) {
	c.ProfilesDirPath = path
}

func (c *Config) SetOllama(host, model string) {
	c.Ollama.Host = host
	c.Ollama.Model = model
}
