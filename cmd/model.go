package cmd

import (
	"github.com/robertoseba/gennie/internal/cache"
	models "github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewModelCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdModel := &cobra.Command{
		Use:   "model",
		Short: "Configures the model to use.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			configModel(c)
		},
	}

	return cmdModel
}

func configModel(c *cache.Cache) {
	model := output.MenuModel(models.ListModels(), models.ModelEnum(c.Model))

	if model == "" {
		return
	}

	if string(model) == c.Model {
		return
	}

	c.SetModel(string(model))

	if err := c.Save(); err != nil {
		ExitWithError(err)
	}
}
