package cmd

import (
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewRootCmd(version string, storage common.IStorage, p *output.Printer, h httpclient.IHttpClient) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gennie",
		Short: "Gennie is a cli assistant with multiple models and profile support.",
	}

	rootCmd.AddCommand(NewModelCmd(storage, p))
	rootCmd.AddCommand(NewProfilesCmd(storage, p))
	rootCmd.AddCommand(NewStatusCmd(storage, p))
	rootCmd.AddCommand(NewAskCmd(storage, p, h))
	rootCmd.AddCommand(NewExportCmd(storage, p))
	rootCmd.AddCommand(NewClearCmd(storage, p))
	rootCmd.AddCommand(NewVersionCmd(version))
	rootCmd.AddCommand(NewConfigCmd(storage, p))
	return rootCmd
}
