package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/models/profile"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdProfiles)
	cmdProfiles.AddCommand(cmdRefreshProfiles)
	cmdProfiles.AddCommand(cmdListProfiles)
	cmdProfiles.AddCommand(cmdConfigProfile)
}

var cmdProfiles = &cobra.Command{
	Use:   "profile",
	Short: "Profile management",
}

var cmdListProfiles = &cobra.Command{
	Use:   "list",
	Short: "List available profiles slug for use with --profile flag when asking questions",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()

		if len(c.Cache.ProfileFilenames) == 0 {
			refreshProfiles(c)
			c.Cache.Save()
		}

		if c.Cache.Profile != nil {
			c.Printer.PrintLine(output.Yellow)
			c.Printer.Print(fmt.Sprintf("Current Profile: %s (%s)", c.Cache.Profile.Name, c.Cache.Profile.Slug), output.Cyan)
		}

		c.Printer.PrintLine(output.Yellow)
		c.Printer.Print("Available Profiles: ", output.Cyan)
		for slug := range c.Cache.ProfileFilenames {
			c.Printer.Print(slug, output.Gray)
		}
		c.Printer.PrintLine(output.Yellow)
	},
}

var cmdRefreshProfiles = &cobra.Command{
	Use:   "refresh",
	Short: "Rescan the profiles folder and update the cache with available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()
		refreshProfiles(c)
		c.Cache.Save()
	},
}

func refreshProfiles(c *Container) {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		ExitWithError(err)
	}

	c.Cache.ProfileFilenames = make(map[string]string, len(profiles))
	for i := range profiles {
		c.Cache.ProfileFilenames[profiles[i].Slug] = profiles[i].Filename
	}
}

var cmdConfigProfile = &cobra.Command{
	Use:   "config",
	Short: "Configures which profile to use.",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()
		configProfile(c)
	},
}

func configProfile(c *Container) {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		ExitWithError(err)
	}
	profileSlug := output.MenuProfile(profiles)

	if profileSlug == "" {
		return
	}

	c.Cache.SetProfile(profiles[profileSlug])

	c.Cache.ProfileFilenames = make(map[string]string, len(profiles))
	for i := range profiles {
		c.Cache.ProfileFilenames[profiles[i].Slug] = profiles[i].Filename
	}

	if err := c.Cache.Save(); err != nil {
		ExitWithError(err)
	}
}
