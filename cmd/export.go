package cmd

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewExportCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdExport := &cobra.Command{
		Use:   "export",
		Short: "Export the chat history to a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p.Print(fmt.Sprintf("Exporting chat history to file %s", args[0]), output.Cyan)
			err := exportChatHistory(c, args[0])
			if err != nil {
				ExitWithError(err)
			}
		},
	}

	return cmdExport
}

func exportChatHistory(c *cache.Cache, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	for _, chat := range c.ChatHistory.Responses {
		var err error
		_, err = f.WriteString(fmt.Sprintf("Question: %s\n", chat.GetQuestion()))
		if err != nil {
			return err
		}

		_, err = f.WriteString(fmt.Sprintf("Answer: %s\n", chat.GetAnswer()))
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
