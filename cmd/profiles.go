package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/models/profile"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewProfilesCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdProfiles := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
	}

	cmdListProfiles := &cobra.Command{
		Use:   "list",
		Short: "List available profiles slug for use with --profile flag when asking questions",
		Run: func(cmd *cobra.Command, args []string) {

			if len(c.ProfileFilenames) == 0 {
				refreshProfiles(c)
				c.Save()
			}

			if c.Profile != nil {
				p.PrintLine(output.Yellow)
				p.Print(fmt.Sprintf("Current Profile: %s (%s)", c.Profile.Name, c.Profile.Slug), output.Cyan)
			}

			p.PrintLine(output.Yellow)
			p.Print("Available Profiles: ", output.Cyan)
			for slug := range c.ProfileFilenames {
				p.Print(slug, output.Gray)
			}
			p.PrintLine(output.Yellow)
		},
	}

	cmdRefreshProfiles := &cobra.Command{
		Use:   "refresh",
		Short: "Rescan the profiles folder and update the cache with available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			refreshProfiles(c)
			c.Save()
		},
	}

	cmdConfigProfile := &cobra.Command{
		Use:   "config",
		Short: "Configures which profile to use.",
		Run: func(cmd *cobra.Command, args []string) {
			configProfile(c)
		},
	}

	cmdProfiles.AddCommand(cmdRefreshProfiles)
	cmdProfiles.AddCommand(cmdListProfiles)
	cmdProfiles.AddCommand(cmdConfigProfile)

	return cmdProfiles
}

func refreshProfiles(c *cache.Cache) {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		ExitWithError(err)
	}

	c.ProfileFilenames = make(map[string]string, len(profiles))
	for i := range profiles {
		c.ProfileFilenames[profiles[i].Slug] = profiles[i].Filename
	}
}

func configProfile(c *cache.Cache) {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		ExitWithError(err)
	}
	profileSlug := output.MenuProfile(profiles)

	if profileSlug == "" {
		return
	}

	c.SetProfile(profiles[profileSlug])

	c.ProfileFilenames = make(map[string]string, len(profiles))
	for i := range profiles {
		c.ProfileFilenames[profiles[i].Slug] = profiles[i].Filename
	}

	if err := c.Save(); err != nil {
		ExitWithError(err)
	}
}
