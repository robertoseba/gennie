package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewExportCmd(exportConversationCmd *usecases.ExportConversationService, p *output.Printer) *cobra.Command {
	cmdExport := &cobra.Command{
		Use:   "export [filename]",
		Short: "Export the chat history to a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Print(fmt.Sprintf("Exporting chat history to file %s", args[0]), output.Cyan)
			err := exportConversationCmd.Execute(args[0])
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmdExport
}
