package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewClearCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clears all the conversation and preferences from cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Print(fmt.Sprintf("Cache file %s will be removed", c.FilePath), output.Gray)
			p.Print("Press Enter to continue... Or Ctrl+C to cancel", output.Red)
			fmt.Scanln()

			err := c.Clear()
			if err != nil {
				ExitWithError(err)
			}

			p.Print("Cache cleared", output.Red)
			return nil
		},
	}

	return clearCmd
}
