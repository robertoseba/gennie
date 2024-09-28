package cmd

import (
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(c *cache.Cache, p *output.Printer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}

	rootCmd.AddCommand(NewConfigCmd(c, p))
	rootCmd.AddCommand(NewProfilesCmd(c, p))
	rootCmd.AddCommand(NewAskCmd(c, p))

	return rootCmd
}
