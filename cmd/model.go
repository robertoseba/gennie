package cmd

import (
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/usecases"
	output "github.com/robertoseba/gennie/internal/output"
	cobra "github.com/spf13/cobra"
)

func NewModelCmd(selectModelCmd *usecases.SelectModelService, p *output.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Configures the model to use and list slugs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			modelList := selectModelCmd.ListAll()
			modelSelected := output.MenuModel(modelList, models.DefaultModel)

			err := selectModelCmd.SetAsActive(modelSelected)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(newModelSlugsCmd(selectModelCmd, p))
	return cmd
}

func newModelSlugsCmd(selectModelCmd *usecases.SelectModelService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "slugs",
		Short: "List available models slugs for use with --model(-m=) flag when asking questions",
		RunE: func(cmd *cobra.Command, args []string) error {
			p.PrintLine(output.Yellow)
			p.Print("Available Models: ", output.Cyan)

			modelList := selectModelCmd.ListAll()
			for model := range modelList {
				p.Print(model.Slug(), output.Gray)
			}
			p.PrintLine(output.Yellow)

			return nil
		},
	}
}
