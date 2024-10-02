package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of the Gennie",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version)
		},
	}
}
