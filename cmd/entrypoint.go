package cmd

import (
	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(printer *output.Printer,
	version string,
	config *config.Config,
	configRepo config.IConfigRepository,
	askCmd *usecases.GetAnswerService,
	selectModelCmd *usecases.SelectModelService,
	selectProfileCmd *usecases.SelectProfileService,
	exportConversationCmd *usecases.ExportConversationService) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Gennie version: {{.Version}}")

	rootCmd.AddCommand(NewModelCmd(selectModelCmd, printer))
	rootCmd.AddCommand(NewProfilesCmd(selectProfileCmd, printer))
	rootCmd.AddCommand(NewAskCmd(askCmd, printer))
	rootCmd.AddCommand(NewExportCmd(exportConversationCmd, printer))
	rootCmd.AddCommand(NewConfigCmd(configRepo, printer))
	rootCmd.AddCommand(NewStatusCmd(configRepo, printer))
	return rootCmd
}
