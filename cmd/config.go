package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdConfig)
	cmdConfig.AddCommand(cmdModelConfig)
	cmdConfig.AddCommand(cmdConfigProfile)
}

var cmdConfig = &cobra.Command{
	Use:   "config",
	Short: "Configurations for gennie",
	Long:  `Use it to configure ginnie. You can change models, profiles and more!`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("starting full configuration...")
	},
}

var cmdModelConfig = &cobra.Command{
	Use:   "model",
	Short: "Configures which model to use.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting model configuration...")
	},
}

var cmdConfigProfile = &cobra.Command{
	Use:   "profile",
	Short: "Configures which profile to use.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting profile configuration...")
	},
}
