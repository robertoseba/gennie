package cmd

import (
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/output"
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

			for slug, _ := range availableProfiles {
				p.Print(slug, output.Gray)
			}
			p.PrintLine(output.Yellow)
		},
	}

	cmdRefreshProfiles := &cobra.Command{
		Use:   "refresh",
		Short: "Rescan the profiles folder and update the cache with available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			// config := persistence.GetConfig()
			// scan -> config.ProfilesPath
			// persistence.SetProfileFiles(profileFiles)
			panic("To be implemented")
			p.Print("Profiles refreshed...", output.Cyan)
		},
	}

	cmdProfiles.AddCommand(cmdRefreshProfiles)
	cmdProfiles.AddCommand(cmdListProfiles)

	return cmdProfiles
}
