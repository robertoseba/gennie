package common

type Config struct {
	ProfilesPath    string
	OpenAiApiKey    string
	AnthropicApiKey string
	MaritacaApiKey  string
	StyledTerminal  bool
}

func NewConfig() Config {
	return Config{
		ProfilesPath:    "",
		OpenAiApiKey:    "",
		AnthropicApiKey: "",
		MaritacaApiKey:  "",
		StyledTerminal:  true,
	}
}
