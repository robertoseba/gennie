package cmd

import (
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
	setupSubCommands(command, container, printer)

	if container.GetConfig().IsNew() {
		command.SetArgs([]string{"config"})
	}

	err := command.Execute()
	if err != nil {
		command.PrintErrf("Error executing command: %v", err)
		os.Exit(1)
	}
}

func newRootCmd(version string, stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "gennie",
		Short:             "Gennie is a cli assistant with multiple models and profile support.",
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Gennie version: {{.Version}}")

	return rootCmd
}

func setupSubCommands(c *cobra.Command, container *container.Container, printer *output.Printer) {
	subcmds := []*cobra.Command{
		NewModelCmd(container.GetSelectModelService(), printer),
		NewProfilesCmd(container.GetSelectProfileService(), printer),
		NewAskCmd(container.GetCompleteService(), printer),
		NewExportCmd(container.GetExportConversationService(), printer),
		NewConfigCmd(container.GetConfigRepository(), printer),
		NewStatusCmd(container.GetConfigRepository(), printer),
	}
	c.AddCommand(subcmds...)
}
