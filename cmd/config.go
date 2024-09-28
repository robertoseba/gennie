package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
	models "github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewConfigCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdConfig := &cobra.Command{
		Use:   "config",
		Short: "Configurations for gennie",
		Long:  `Use it to configure ginnie. You can change models, profiles and more!`,
		Run: func(cmd *cobra.Command, args []string) {
			configModel(c)
			configProfile(c) // todo: where to better put this function
		},
	}

	cmdShowConfig := &cobra.Command{
		Use:   "show",
		Short: "Shows the current configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			if c.Model == "" || c.Profile == nil {
				ExitWithError(fmt.Errorf("Gennie hasn't been configured yet. Please run gennie config first."))
			}

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Model: %s ", models.ModelEnum(c.Model)), output.Cyan)
			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profile name: %s", c.Profile.Name), output.Gray)
			p.Print(fmt.Sprintf("Profile description: %s", c.Profile.Data), output.Gray)
			p.PrintLine(output.Yellow)
		},
	}

	cmdModelConfig := &cobra.Command{
		Use:   "model",
		Short: "Configures which model to use.",
		Run: func(cmd *cobra.Command, args []string) {
			configModel(c)
		},
	}

	cmdConfig.AddCommand(cmdModelConfig)
	cmdConfig.AddCommand(cmdShowConfig)

	return cmdConfig
}

func configModel(c *cache.Cache) {
	model := output.MenuModel()

	if model == "" {
		return
	}

	c.SetModel(string(model))

	if err := c.Save(); err != nil {
		ExitWithError(err)
	}
}
