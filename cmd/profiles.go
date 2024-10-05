package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/profile"
	"github.com/spf13/cobra"
)

func NewProfilesCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdProfiles := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
		Run: func(cmd *cobra.Command, args []string) {
			configProfile(c, p)
		},
	}

	cmdListProfiles := &cobra.Command{
		Use:   "slugs",
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
			p.Print("Profiles refreshed...", output.Cyan)
		},
	}

	cmdProfiles.AddCommand(cmdRefreshProfiles)
	cmdProfiles.AddCommand(cmdListProfiles)

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

func configProfile(c *cache.Cache, p *output.Printer) {
	if c.Profile == nil {
		c.Profile = profile.DefaultProfile()
	}

	profiles, err := profile.LoadProfiles()
	if err != nil {
		p.Print(fmt.Sprintf("Could not load profiles: %s\n", err.Error()), output.Red)
		profiles = make(map[string]*profile.Profile, 1)
		profiles[c.Profile.Slug] = c.Profile
	}

	profileSlug := output.MenuProfile(profiles, c.Profile.Slug)

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
