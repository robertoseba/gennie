package profile

import (
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/robertoseba/gennie/internal/core/config"
)

type ProfileRepositoryInterface interface {
	ListAll() map[string]string
	FindBySlug(slug string) (*Profile, error)
	RefreshProfiles() error
}

type ProfileRepository struct {
	profilesInfo map[string]ProfileInfo
	config       *config.Config
}

func NewProfileRepository(config *config.Config, profilesInfo map[string]ProfileInfo) *ProfileRepository {
	return &ProfileRepository{
		profilesInfo: profilesInfo,
		config:       config,
	}
}

// Returns a map with all the profiles [slug]name
func (pr *ProfileRepository) ListAll() map[string]string {
	allProfiles := make(map[string]string, len(pr.profilesInfo))

	for slug := range pr.profilesInfo {
		allProfiles[slug] = pr.profilesInfo[slug].Name
	}

	return allProfiles
}

func (pr *ProfileRepository) FindBySlug(slug string) (*Profile, error) {
	if slug == DefaultProfileSlug {
		return DefaultProfile(), nil
	}

	profileInfo, ok := pr.profilesInfo[slug]
	if !ok {
		return nil, ErrNoProfileSlug
	}

	return pr.loadProfileFile(profileInfo.Filepath)
}

func (pr *ProfileRepository) RefreshProfiles() error {
	prDir := pr.config.ProfilesDirPath

	files, err := os.ReadDir(prDir)
	if err != nil {
		return ErrProfileDirNotExist
	}

	profilesInfo := make(map[string]ProfileInfo, len(files))

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if strings.HasSuffix(filename, profileFileExtension) {
				//TODO: load profiles in parallel
				loadedProfile, err := pr.loadProfileFile(path.Join(prDir, file.Name()))
				if err != nil {
					return err
				}

				if loadedProfile.Slug == "" {
					continue
				}

				profilesInfo[loadedProfile.Slug] = ProfileInfo{
					Slug:     loadedProfile.Slug,
					Name:     loadedProfile.Name,
					Filepath: path.Join(prDir, file.Name()),
				}

			}
		}
	}

	pr.profilesInfo = profilesInfo

	return nil
}

func (pr *ProfileRepository) loadProfileFile(filePath string) (*Profile, error) {
	profile := &Profile{}
	_, err := toml.DecodeFile(filePath, profile)

	if err != nil {
		return nil, ErrLoadingToml
	}

	return profile, nil
}
