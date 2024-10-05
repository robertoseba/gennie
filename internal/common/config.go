package common

import "github.com/robertoseba/gennie/internal/profile"

type Config struct {
	CurrModelSlug   string
	CurrProfile     profile.Profile
	ProfilesPath    string
	OpenAiApiKey    string
	AnthropicApiKey string
	MaritacaApiKey  string
	StyledTerminal  bool
}

func NewConfig() Config {
	return Config{
		CurrModelSlug:   "default",
		CurrProfile:     *profile.DefaultProfile(),
		ProfilesPath:    "",
		OpenAiApiKey:    "",
		AnthropicApiKey: "",
		MaritacaApiKey:  "",
		StyledTerminal:  true,
	}
}

// func (c *Config) GobEncode() ([]byte, error) {
// 	buffer := new(bytes.Buffer)
// 	encoder := gob.NewEncoder(buffer)
// 	err := encoder.Encode(c)
// 	return buffer.Bytes(), err
// }

// func (c *Config) GobDecode(data []byte) error {
// 	buffer := bytes.NewBuffer(data)
// 	decoder := gob.NewDecoder(buffer)
// 	return decoder.Decode(c)
// }
