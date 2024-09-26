package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdAsk)
	cmdAsk.Flags().BoolVarP(&isFollowUp, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&askFile, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&model, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profile, "profile", "p", "", "specifies the profile to use.")

}

var isFollowUp bool
var askFile string
var model string
var profile string

var cmdAsk = &cobra.Command{
	Use:   "ask [question for the llm model]",
	Short: "You can ask anything here",
	Long:  `The question that will be sent to the llm. If your question contains special characters, please use quotes.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if model != "" && !slices.Contains(models.ListModels(), models.ModelEnum(model)) {
			return fmt.Errorf("Model %s not supported. Please use one of the following:\n%s\n", model, strings.Join(models.ListModelsSlug(), ", "))
		}

		fmt.Println("Question: " + strings.Join(args, " "))
		fmt.Println("is followup: ", isFollowUp)
		fmt.Println("Appending file: ", askFile)
		fmt.Println("Model: ", model)
		fmt.Println("Profile: ", profile)
		return nil
	},
}
