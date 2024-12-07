package repositories

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/robertoseba/gennie/internal/core/profile"
)

var ErrNoProfilesDir = fmt.Errorf("no profiles found. Please add profiles to the profiles folder.")

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
// If dir does not exist, returns the default profile only and the error ErrNoProfilesDir
func (pr *ProfileRepository) ListAll() (map[string]*profile.Profile, error) {
	profiles, err := pr.scanProfiles()

	if err != nil {
		if errors.Is(err, ErrNoProfilesDir) {
			return map[string]*profile.Profile{profile.DefaultProfileSlug: profile.DefaultProfile()}, err
		}
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
		return nil, ErrNoProfilesDir
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

// Uses the os.UserConfigDir to setup the default profiles directory
func DefaultProfileDir() string {
	const gennieConfigDir = "gennie"
	const gennieProfilesDir = "profiles"

	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return path.Join(configDir, gennieConfigDir, gennieProfilesDir)
}
