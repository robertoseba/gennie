package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewClearCmd(storage common.IStorage, p *output.Printer) *cobra.Command {

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clears all the conversation and preferences from cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Print("All your cached data will be erased!", output.Gray)
			p.Print("Press Enter to continue... Or Ctrl+C to cancel", output.Red)
			fmt.Scanln()

			storage.Clear()

			p.Print("Cache cleared", output.Red)
			return nil
		},
	}

	return clearCmd
}
