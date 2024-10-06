package profile

import (
	"os"
	"path"
)

type Profile struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Author string `json:"author"`
	Data   string `json:"data"`
}

type ProfileInfo struct {
	Slug     string
	Name     string
	Filepath string
}

func DefaultProfile() *Profile {
	return &Profile{
		Name:   "Default assistant",
		Author: "gennie",
		Slug:   "default",
		Data:   "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necessary.",
	}
}

func DefaultProfilePath() string {
	const gennieConfigDir = "gennie"
	const gennieProfilesDir = "profiles"

	configDir, err := os.UserConfigDir()

	if err != nil {
		return ""
	}

	return path.Join(configDir, gennieConfigDir, gennieProfilesDir)
}
