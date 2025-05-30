package cmd

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/infra/container"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func Run(version string, stdOut io.Writer, stdErr io.Writer) {
	container := container.NewContainer()
	printer := output.NewPrinter(stdOut, stdErr)

	command := newRootCmd(version, stdOut, stdErr)

	slog.SetDefault(slog.New(slog.DiscardHandler))
	if os.Getenv("DEBUG") != "" {
		logFile := setupDebugLogger(path.Join(container.GetConfig().ConversationCacheDir, "gennie.log"))
		defer logFile.Close()
	}

	setupSubCommands(command, container, printer)

	if container.GetConfig().IsNew() {
		command.SetArgs([]string{"config"})
	}

	// Disables terminal output if not running in a terminal
	command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if !term.IsTerminal(int(os.Stdout.Fd())) {
			err := cmd.Flags().Set("terminal", "false")
			if err != nil {
				cmd.PrintErrf("Error setting terminal flag: %v\n", err)
				os.Exit(1)
			}
		}
	}

	err := command.ExecuteContext(context.Background())
	if err != nil {
		command.PrintErrf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd(version string, stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	var isTerminalFlag bool

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
	rootCmd.PersistentFlags().BoolVarP(&isTerminalFlag, "terminal", "t", true, "controls if output if interactive terminal or if should just output plain data. Can be used for piping")

	return rootCmd
}

func setupSubCommands(c *cobra.Command, container *container.Container, printer *output.Printer) {
	subcmds := []*cobra.Command{
		NewModelCmd(container.GetSelectModelService(), printer),
		NewProfilesCmd(container.GetSelectProfileService(), printer),
		NewAskCmd(container.GetCompleteService(), printer),
		NewConfigCmd(container.GetConfigRepository(), printer),
		NewStatusCmd(container.GetConfigRepository(), printer),
		NewConversationCmd(container.GetExportConversationService(), printer),
	}
	c.AddCommand(subcmds...)
}

func setupDebugLogger(logFilename string) *os.File {
	logFile, err := os.OpenFile(logFilename, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening log file %s: %v", logFilename, err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)
	return logFile
}
