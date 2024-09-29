package cmd

import (
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(c *cache.Cache, p *output.Printer, h httpclient.IHttpClient) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}

	rootCmd.AddCommand(NewModelCmd(c, p))
	rootCmd.AddCommand(NewProfilesCmd(c, p))
	rootCmd.AddCommand(NewAskCmd(c, p, h))
	rootCmd.AddCommand(NewStatusCmd(c, p))
	rootCmd.AddCommand(NewExportCmd(c, p))

	return rootCmd
}
