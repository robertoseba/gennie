package cmd

import (
	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewProfilesCmd(selectProfileCmd *usecases.SelectProfileService, p *output.Printer) *cobra.Command {
	cmdProfiles := &cobra.Command{
		Use:   "profile",
		Short: "Configures the profile to use and list slugs",
		RunE: func(cmd *cobra.Command, args []string) error {
			availableProfiles, err := selectProfileCmd.ListAll()
			if err != nil {
				if len(availableProfiles) == 0 {
					return err
				}
				p.Print(err.Error(), output.Red)
			}

			menuMap := make(map[string]string, len(availableProfiles))
			for slug, profile := range availableProfiles {
				menuMap[slug] = profile.Name
			}

			selectedProfileSlug := output.MenuProfile(menuMap, "default")

			// When esc is pressed in the menu
			if selectedProfileSlug == "" {
				return nil
			}

			err = selectProfileCmd.SetAsActive(availableProfiles[selectedProfileSlug])
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmdProfiles.AddCommand(newCmdListProfiles(selectProfileCmd, p))
	return cmdProfiles
}

func newCmdListProfiles(selectProfileCmd *usecases.SelectProfileService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "slugs",
		Short: "List available profiles slugs for use with --profile flag when asking questions",
		Long:  "List available profiles slugs for use with --profile(-p=) flag when asking questions. Profile slugs are derived from the filename. Ie: \"my_profile.profile.toml\" will have the slug \"my_profile\".",
		RunE: func(cmd *cobra.Command, args []string) error {
			p.PrintLine(output.Yellow)

			availableProfiles, err := selectProfileCmd.ListAll()
			if err != nil {
				return err
			}

			if len(availableProfiles) == 0 {
				p.Print("No profiles found. Please add profiles to the profiles folder.", output.Red)
				return err
			}

			p.Print("Available Profiles: ", output.Cyan)

			for slug := range availableProfiles {
				p.Print(slug, output.Gray)
			}
			p.PrintLine(output.Yellow)
			return nil
		},
	}
}
