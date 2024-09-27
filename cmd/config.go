package cmd

import (
	"fmt"
	"strconv"

	models "github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdConfig)
	cmdConfig.AddCommand(cmdModelConfig)
	cmdConfig.AddCommand(cmdConfigProfile)
	cmdConfig.AddCommand(cmdShowConfig)
}

var cmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Configurations for gennie",
	Long:  `Use it to configure ginnie. You can change models, profiles and more!`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()
		configModel(c)
		configProfile(c)
	},
}

var cmdShowConfig = &cobra.Command{
	Use:   "show",
	Short: "Shows the current configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()

		c.Printer.PrintLine(output.Yellow)
		c.Printer.Print(fmt.Sprintf("Model: %s ", models.ModelEnum(c.Cache.Model)), output.Cyan)
		c.Printer.PrintLine(output.Yellow)
		c.Printer.Print(fmt.Sprintf("Profile name: %s", c.Cache.Profile.Name), output.Gray)
		c.Printer.Print(fmt.Sprintf("Profile description: %s", c.Cache.Profile.Data), output.Gray)
		c.Printer.PrintLine(output.Yellow)
	},
}

var cmdModelConfig = &cobra.Command{
	Use:   "model",
	Short: "Configures which model to use.",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()
		configModel(c)
	},
}

var cmdConfigProfile = &cobra.Command{
	Use:   "profile",
	Short: "Configures which profile to use.",
	Run: func(cmd *cobra.Command, args []string) {
		c := setUp()
		configProfile(c)
	},
}

func configModel(c *Container) {
	model := output.MenuModel()
	c.Cache.SetModel(string(model))

	if err := c.Cache.Save(); err != nil {
		ExitWithError(err)
	}
}

func configProfile(c *Container) {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		ExitWithError(err)
	}
	profileIdx := output.MenuProfile(&profiles)

	idx, err := strconv.Atoi(profileIdx)

	if err != nil {
		ExitWithError(err)
	}

	c.Cache.SetProfile(&profiles[idx])

	if err := c.Cache.Save(); err != nil {
		ExitWithError(err)
	}
}
