package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/robertoseba/gennie/internal/infra/container"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func Run(version string, stdOut io.Writer, stdErr io.Writer) {
	container := container.NewContainer()
	printer := output.NewPrinter(stdOut, stdErr)

	command := newRootCmd(version, stdOut, stdErr)
	subcmds := []*cobra.Command{
		NewModelCmd(container.GetSelectModelService(), printer),
		NewProfilesCmd(container.GetSelectProfileService(), printer),
		NewAskCmd(container.GetCompleteService(), printer),
		NewExportCmd(container.GetExportConversationService(), printer),
		NewConfigCmd(container.GetConfigRepository(), printer),
		NewStatusCmd(container.GetConfigRepository(), printer),
	}
	addSubCommands(command, subcmds)

	if container.GetConfig().IsNew() {
		command.SetArgs([]string{"config"})
	}

	err := command.Execute()
	if err != nil {
		stdErr := command.OutOrStderr()
		fmt.Fprintf(stdErr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd(version string, stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Gennie version: {{.Version}}")

	return rootCmd
}

func addSubCommands(rootCmd *cobra.Command, commands []*cobra.Command) {
	rootCmd.AddCommand(commands...)
}
