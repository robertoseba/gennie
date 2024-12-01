package repositories

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type ProfileRepository struct {
	profilesDir   string
	profileSuffix string
}

func NewProfileRepository(profilesDir string) *ProfileRepository {
	return &ProfileRepository{
		profilesDir:   profilesDir,
		profileSuffix: ".profile.toml",
	}
}

// Lists all profiles found in the profiles directory, plus a default profile created in the application
func (pr *ProfileRepository) ListAll() (map[string]*profile.Profile, error) {
	profiles, err := pr.scanProfiles()

	if err != nil {
		return nil, err
	}

	profiles[profile.DefaultProfileSlug] = profile.DefaultProfile()
	return profiles, nil
}

// Loads a profile from a toml file named as the slug. Ie: "test.profile.toml" for slug "test"
func (pr *ProfileRepository) FindBySlug(slug string) (*profile.Profile, error) {
	if slug == profile.DefaultProfileSlug {
		return profile.DefaultProfile(), nil
	}

	file := path.Join(pr.profilesDir, slug+pr.profileSuffix)

	return pr.loadProfileFile(file)
}

// Scans the profiles directory and loads all profiles found
func (pr *ProfileRepository) scanProfiles() (map[string]*profile.Profile, error) {
	prDir := pr.profilesDir

	files, err := os.ReadDir(prDir)
	if err != nil {
		return nil, fmt.Errorf("error reading profiles directory: %w", err)
	}

	profiles := make(map[string]*profile.Profile, len(files))

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if strings.HasSuffix(filename, pr.profileSuffix) {
				//TODO: load profiles in parallel
				loadedProfile, err := pr.loadProfileFile(path.Join(prDir, file.Name()))
				if err != nil {
					return nil, err
				}

				if loadedProfile.Name == "" {
					continue
				}

				loadedProfile.Slug = strings.TrimSuffix(filename, pr.profileSuffix)

				profiles[loadedProfile.Slug] = loadedProfile
			}
		}
	}

	return profiles, nil
}

func (pr *ProfileRepository) loadProfileFile(filePath string) (*profile.Profile, error) {
	profile := &profile.Profile{}
	_, err := toml.DecodeFile(filePath, profile)

	if err != nil {
		return nil, fmt.Errorf("error loading toml file: %w", err)
	}

	profile.Slug = strings.TrimSuffix(path.Base(filePath), pr.profileSuffix)

	return profile, nil
}
