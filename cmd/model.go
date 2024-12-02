package cmd

import (
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/usecases"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewModelCmd(selectModelCmd usecases.SelectModelService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "model",
		Short: "Configures the model to use.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			modelList := selectModelCmd.ListAll()
			modelSelected := output.MenuModel(modelList, models.DefaultModel)
			err := selectModelCmd.SetAsActive(modelSelected)

			//TODO: refactor this to use cmd output e not printer anymore
			ExitWithError(err)
		},
	}
}
