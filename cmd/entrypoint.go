package cmd

import (
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string, persistence common.IPersistence, p *output.Printer, h httpclient.IHttpClient) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}

	rootCmd.AddCommand(NewModelCmd(persistence, p))
	rootCmd.AddCommand(NewProfilesCmd(persistence, p))
	rootCmd.AddCommand(NewStatusCmd(persistence, p))
	// rootCmd.AddCommand(NewAskCmd(persistence, p, h))
	// rootCmd.AddCommand(NewExportCmd(persistence, p))
	// rootCmd.AddCommand(NewClearCmd(persistence, p))
	rootCmd.AddCommand(NewVersionCmd(version))
	return rootCmd
}
