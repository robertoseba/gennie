package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "gennie"}

func Execute() {
	rootCmd.Execute()
}
