package profile

const DefaultProfileSlug = "default"

type Profile struct {
	Name       string               `toml:"name"`
	Slug       string               `toml:"_"`
	Author     string               `toml:"author"`
	Data       string               `toml:"data"`
	McpServers map[string]mcpServer `toml:"mcpServers"`
}

func DefaultProfile() *Profile {
	return &Profile{
		Name:   "Default assistant",
		Author: "gennie",
		Slug:   DefaultProfileSlug,
		Data:   "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necessary.",
	}
}

type mcpServer struct {
	Command          string   `toml:"command"`
	Args             []string `toml:"args"`
	RequiresApproval bool     `toml:"requires_approval"`
	Envs             []string `toml:"envs"`
	AllowedTools     []string `toml:"allowed_tools"`
}
