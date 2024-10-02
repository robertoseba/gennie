package profile

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
)

type Profile struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Author   string `json:"author"`
	Data     string `json:"data"`
	Filename string
}

func LoadProfiles() (map[string]*Profile, error) {
	const gennieConfigDir = "gennie"
	const gennieProfilesDir = "profiles"

	profilesPath := os.Getenv("GINNIE_PROFILES_PATH")
	if profilesPath == "" {
		configDir, err := os.UserConfigDir()
		profilesPath = path.Join(configDir, gennieConfigDir, gennieProfilesDir)

		_, errStats := os.Stat(profilesPath)

		if err != nil || os.IsNotExist(errStats) {
			curr, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			profilesPath = path.Join(curr, gennieProfilesDir)
		}
	}

	if _, err := os.Stat(profilesPath); os.IsNotExist(err) {
		return nil, ErrNoFolder{Path: profilesPath}
	}

	profileFiles, err := loadFileNames(profilesPath)

	if err != nil {
		return nil, err
	}
	if len(profileFiles) == 0 {
		return nil, ErrNoProfiles{profilesPath}
	}

	profiles := make(map[string]*Profile, len(profileFiles)+1)

	defaultProfile := LoadDefaultProfile()
	profiles[defaultProfile.Slug] = defaultProfile

	for _, profileFile := range profileFiles {
		profile, err := LoadProfileFromFile(profileFile)
		profile.Filename = profileFile
		if err != nil {
			return nil, err
		}
		profiles[profile.Slug] = profile
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

func LoadDefaultProfile() *Profile {
	return &Profile{
		Name:   "Default assistant",
		Author: "gennie",
		Slug:   "default",
		Data:   "You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necesary.",
	}
}
