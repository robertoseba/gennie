package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var isFollowUp bool
	var askFile string

	var cmdAsk = &cobra.Command{
		Use:   "ask [question for the llm model]",
		Short: "You can ask anything here",
		Long:  `The question that will be sent to the llm. If your question contains special characters, please use quotes.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Question: " + strings.Join(args, " "))
			fmt.Println("is followup: ", isFollowUp)
			fmt.Println("Appending file: ", askFile)
		},
	}

	var cmdConfig = &cobra.Command{
		Use:   "config [what would you like to configure]",
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

	cmdAsk.Flags().BoolVarP(&isFollowUp, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&askFile, "append", "a", "", "appends the content of a file to the question.")

	var rootCmd = &cobra.Command{Use: "gennie"}
	rootCmd.AddCommand(cmdConfig, cmdAsk)
	cmdConfig.AddCommand(cmdModelConfig)
	cmdConfig.AddCommand(cmdConfigProfile)
	rootCmd.Execute()
}
