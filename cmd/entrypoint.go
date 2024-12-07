package cmd

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/infra/container"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func Run(version string) {
	container := container.NewContainer()
	printer := output.NewPrinter(nil, nil)

	command := newRootCmd(version, printer, container)
	if container.GetConfig().IsNew() {
		command.SetArgs([]string{"config"})
	}

	command.SetOut(os.Stdout)
	command.SetErr(os.Stderr)
	err := command.Execute()

	if err != nil {
		stdErr := command.OutOrStderr()
		fmt.Fprintf(stdErr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
func newRootCmd(version string, printer *output.Printer, container *container.Container) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Gennie version: {{.Version}}")

	rootCmd.AddCommand(NewModelCmd(container.GetSelectModelService(), printer))
	rootCmd.AddCommand(NewProfilesCmd(container.GetSelectProfileService(), printer))
	rootCmd.AddCommand(NewAskCmd(container.GetCompleteService(), printer))
	rootCmd.AddCommand(NewExportCmd(container.GetExportConversationService(), printer))
	rootCmd.AddCommand(NewConfigCmd(container.GetConfigRepository(), printer))
	rootCmd.AddCommand(NewStatusCmd(container.GetConfigRepository(), printer))

	return rootCmd
}
