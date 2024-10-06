package cmd

import (
	"github.com/robertoseba/gennie/internal/common"
	models "github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewModelCmd(storage common.IStorage, p *output.Printer) *cobra.Command {

	cmdModel := &cobra.Command{
		Use:   "model",
		Short: "Configures the model to use.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			modelSelected := output.MenuModel(models.ListModels(), models.ModelEnum(storage.GetCurrModelSlug()))
			storage.SetCurrModelSlug(string(modelSelected))
		},
	}

	return cmdModel
}
