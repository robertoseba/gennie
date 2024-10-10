package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(c *CmdUtil) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}
	rootCmd.Version = c.Version
	rootCmd.SetVersionTemplate("Gennie version: {{.Version}}")

	rootCmd.AddCommand(NewModelCmd(c.Storage, c.Printer))
	rootCmd.AddCommand(NewProfilesCmd(c.Storage, c.Printer))
	rootCmd.AddCommand(NewStatusCmd(c.Storage, c.Printer))
	rootCmd.AddCommand(NewAskCmd(c.Storage, c.Printer, c.HttpClient()))
	rootCmd.AddCommand(NewExportCmd(c.Storage, c.Printer))
	rootCmd.AddCommand(NewClearCmd(c.Storage, c.Printer))
	rootCmd.AddCommand(NewConfigCmd(c.Storage, c.Printer))
	return rootCmd
}
