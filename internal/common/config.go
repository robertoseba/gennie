package common

type Config struct {
	ProfilesPath    string
	OpenAiApiKey    string
	AnthropicApiKey string
	MaritacaApiKey  string
	GroqApiKey      string
	StyledTerminal  bool
	OllamaHost      string
	OllamaModel     string
}

func NewConfig() Config {
	return Config{
		ProfilesPath:    "",
		OpenAiApiKey:    "",
		AnthropicApiKey: "",
		MaritacaApiKey:  "",
		GroqApiKey:      "",
		StyledTerminal:  true,
		OllamaHost:      "",
		OllamaModel:     "",
	}
}
