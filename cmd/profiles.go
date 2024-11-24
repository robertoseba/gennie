package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/profile"
	"github.com/spf13/cobra"
)

func NewProfilesCmd(storage common.IStorage, p *output.Printer) *cobra.Command {
	cmdProfiles := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
		Run: func(cmd *cobra.Command, args []string) {
			availableProfiles := storage.GetCachedProfiles()

			if len(availableProfiles) == 0 {
				p.Print("No profiles found. Please add profiles to the profiles folder.", output.Red)
				return
			}

			menuMap := make(map[string]string, len(availableProfiles))
			for slug, profile := range availableProfiles {
				menuMap[slug] = profile.Name
			}

			selectedProfileSlug := output.MenuProfile(menuMap, storage.GetCurrProfile().Slug)

			if selectedProfileSlug == "" {
				return
			}

			profile, err := storage.LoadProfileData(selectedProfileSlug)
			if err != nil {
				ExitWithError(err)
			}

			storage.SetCurrProfile(*profile)
		},
	}

	cmdListProfiles := &cobra.Command{
		Use:   "slugs",
		Short: "List available profiles slug for use with --profile flag when asking questions",
		Run: func(cmd *cobra.Command, args []string) {
			p.PrintLine(output.Yellow)

			availableProfiles := storage.GetCachedProfiles()

			if len(availableProfiles) == 0 {
				p.Print("No profiles found. Please add profiles to the profiles folder.", output.Red)
				return
			}

			p.Print("Available Profiles: ", output.Cyan)

			for slug := range availableProfiles {
				p.Print(slug, output.Gray)
			}
			p.PrintLine(output.Yellow)
		},
	}

	cmdRefreshProfiles := &cobra.Command{
		Use:   "refresh",
		Short: "Rescan the profiles folder and update the cache with available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			refreshProfiles(storage)
			p.Print("Profiles refreshed...", output.Cyan)
		},
	}

	cmdProfiles.AddCommand(cmdRefreshProfiles)
	cmdProfiles.AddCommand(cmdListProfiles)

	return cmdProfiles
}

func refreshProfiles(storage common.IStorage) {
	config := storage.GetConfig()

	cachedProfiles, err := scanProfilesFolder(config.ProfilesPath)
	if err != nil {
		ExitWithError(err)
	}
	storage.SetCachedProfiles(cachedProfiles)
}

func scanProfilesFolder(dirpath string) (map[string]profile.ProfileInfo, error) {
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	profileInfo := make(map[string]profile.ProfileInfo, len(files))

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if strings.HasSuffix(filename, ".profile.toml") {
				loadedProfile := &profile.Profile{}
				_, err := toml.DecodeFile(path.Join(dirpath, file.Name()), loadedProfile)
				if err != nil {
					return nil, err
				}

				if loadedProfile.Slug == "" {
					continue
				}

				currInfo := profile.ProfileInfo{
					Slug:     loadedProfile.Slug,
					Name:     loadedProfile.Name,
					Filepath: path.Join(dirpath, file.Name()),
				}

				profileInfo[currInfo.Slug] = currInfo
			}
		}
	}

	return profileInfo, nil
}
