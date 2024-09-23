package profile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Profile struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Author string `json:"author"`
	Data   string `json:"data"`
}

func LoadProfiles() ([]Profile, error) {
	const profileDir = "gennie/profiles"

	profilesPath := os.Getenv("GINNIE_PROFILES_PATH")
	if profilesPath == "" {

		configDir, err := os.UserConfigDir()

		if err != nil {
			curr, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			profilesPath = path.Join(curr, profileDir)
		}

		profilesPath = path.Join(configDir, profileDir)
	}

	if _, err := os.Stat(profilesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("No profiles folder in %s", profilesPath)
	}

	profileFiles, err := loadFileNames(profilesPath)

	if err != nil {
		return nil, err
	}
	if len(profileFiles) == 0 {
		return nil, fmt.Errorf("No profiles found in %s", profilesPath)
	}

	profiles := make([]Profile, 0, len(profileFiles)+1)

	profiles = append(profiles, *createDefaultProfile())

	for _, profileFile := range profileFiles {
		profile, err := LoadProfileFromFile(profileFile)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *profile)
	}

	return profiles, nil
}

func loadFileNames(profilesPath string) ([]string, error) {
	profileFiles, err := os.ReadDir(profilesPath)

	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(profileFiles))

	for _, file := range profileFiles {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".profile.json") {
			files = append(files, path.Join(profilesPath, file.Name()))
		}
	}
	return files, nil
}

func LoadProfileFromFile(profilePath string) (*Profile, error) {
	file, err := os.Open(profilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	profile := &Profile{}

	err = json.Unmarshal(data, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func createDefaultProfile() *Profile {
	return &Profile{
		Name:   "Default assistant",
		Author: "gennie",
		Slug:   "default",
		Data:   "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necesary.",
	}
}
