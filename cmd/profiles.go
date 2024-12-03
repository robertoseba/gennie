package cmd

import (
	"errors"

	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/infra/repositories"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewProfilesCmd(selectProfileCmd *usecases.SelectProfileService, p *output.Printer) *cobra.Command {
	cmdProfiles := &cobra.Command{
		Use:   "profile",
		Short: "Configures the profile to use and list slugs",
		Run: func(cmd *cobra.Command, args []string) {
			availableProfiles, err := selectProfileCmd.ListAll()
			if err != nil {
				if !errors.Is(err, repositories.ErrNoProfilesDir) {
					ExitWithError(err)
				}
				p.Print("No profiles found. Please add profiles to the profiles folder.", output.Red)
			}

			menuMap := make(map[string]string, len(availableProfiles))
			for slug, profile := range availableProfiles {
				menuMap[slug] = profile.Name
			}

			selectedProfileSlug := output.MenuProfile(menuMap, "default")

			if selectedProfileSlug == "" {
				return
			}

			err = selectProfileCmd.SetAsActive(availableProfiles[selectedProfileSlug])
			if err != nil {
				ExitWithError(err)
			}
		},
	}

	cmdListProfiles := &cobra.Command{
		Use:   "slugs",
		Short: "List available profiles slugs for use with --profile flag when asking questions",
		Long:  "List available profiles slugs for use with --profile(-p=) flag when asking questions. Profile slugs are derived from the filename. Ie: \"my_profile.profile.toml\" will have the slug \"my_profile\".",
		Run: func(cmd *cobra.Command, args []string) {
			p.PrintLine(output.Yellow)

			availableProfiles, err := selectProfileCmd.ListAll()
			if err != nil {
				ExitWithError(err)
			}

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

	cmdProfiles.AddCommand(cmdListProfiles)
	return cmdProfiles
}
