package profile

import (
	"os"
	"path"
)

const DefaultProfileSlug = "default"

type Profile struct {
	Name   string `toml:"name"`
	Slug   string `toml:"_"`
	Author string `toml:"author"`
	Data   string `toml:"data"`
}

func DefaultProfile() *Profile {
	return &Profile{
		Name:   "Default assistant",
		Author: "gennie",
		Slug:   DefaultProfileSlug,
		Data:   "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necessary.",
	}
}

// TODO: Move this to repository infra
func DefaultProfilePath() string {
	const gennieConfigDir = "gennie"
	const gennieProfilesDir = "profiles"

	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return path.Join(configDir, gennieConfigDir, gennieProfilesDir)
}
