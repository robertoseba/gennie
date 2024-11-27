package cmd

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewExportCmd(storage common.IStorage, p *output.Printer) *cobra.Command {
	cmdExport := &cobra.Command{
		Use:   "export [filename]",
		Short: "Export the chat history to a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p.Print(fmt.Sprintf("Exporting chat history to file %s", args[0]), output.Cyan)
			err := exportChatHistory(storage, args[0])
			if err != nil {
				ExitWithError(err)
			}
		},
	}

	return cmdExport
}

func exportChatHistory(storage common.IStorage, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	chatHistory := storage.GetChatHistory()

	for _, chat := range chatHistory.Responses {
		var err error
		_, err = fmt.Fprintf(f, "## Question: %s\n", chat.GetQuestion())
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(f, "## Answer: %s\n", chat.GetAnswer())
		if err != nil {
			return err
		}

		_, err = f.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}
