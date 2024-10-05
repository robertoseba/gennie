package cmd

import (
	"github.com/robertoseba/gennie/internal/common"
	models "github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewModelCmd(persistence common.IPersistence, p *output.Printer) *cobra.Command {

	cmdModel := &cobra.Command{
		Use:   "model",
		Short: "Configures the model to use.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := persistence.GetConfig()
			modelSelected := output.MenuModel(models.ListModels(), models.ModelEnum(config.CurrModelSlug))
			config.CurrModelSlug = string(modelSelected)
			persistence.SetConfig(config)
		},
	}

	return cmdModel
}
