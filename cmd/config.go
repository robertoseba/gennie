package cmd

import (
	"fmt"

	models "github.com/robertoseba/gennie/internal/models"
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
		if c.Cache.Model == "" || c.Cache.Profile == nil {
			ExitWithError(fmt.Errorf("Gennie hasn't been configured yet. Please run gennie config first."))
		}

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

func configModel(c *Container) {
	model := output.MenuModel()

	if model == "" {
		return
	}

	c.Cache.SetModel(string(model))

	if err := c.Cache.Save(); err != nil {
		ExitWithError(err)
	}
}
