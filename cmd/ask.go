package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdAsk)
	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

}

var isFollowUpFlag bool
var appendFileFlag string
var modelFlag string
var profileFlag string

var cmdAsk = &cobra.Command{
	Use:   "ask [question for the llm model]",
	Short: "You can ask anything here",
	Long:  `The question that will be sent to the llm. If your question contains special characters, please use quotes.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if modelFlag != "" && !slices.Contains(models.ListModels(), models.ModelEnum(modelFlag)) {
			return fmt.Errorf("Model %s not supported. Please use one of the following:\n%s\n", modelFlag, strings.Join(models.ListModelsSlug(), ", "))
		}

		if appendFileFlag != "" {
			return fmt.Errorf("File append not implemented yet.")
		}

		if isFollowUpFlag {
			return fmt.Errorf("Followup not implemented yet.")
		}

		c := setUp()

		input := &InputOptions{
			Question:   strings.Join(args, " "),
			IsFollowUp: isFollowUpFlag,
			AppendFile: appendFileFlag,
			Model:      modelFlag,
			Profile:    profileFlag,
		}

		askModel(c, input)

		return nil
	},
}

type InputOptions struct {
	Question   string
	IsFollowUp bool
	AppendFile string
	Model      string
	Profile    string
}

func askModel(c *Container, input *InputOptions) {
	client := httpclient.NewClient()

	if input.Model != "" {
		ExitWithError(errors.New("Not implemented yet."))
	}

	if input.Profile != "" {
		ExitWithError(errors.New("Not implemented yet."))
	}

	model := models.NewModel(models.ModelEnum(c.Cache.Model), client)
	res, err := model.Ask(input.Question, c.Cache.Profile, nil)

	if err != nil {
		ExitWithError(err)
	}

	c.Printer.PrintLine(output.Yellow)
	c.Printer.PrintAnswer(res.Answer())
	c.Printer.PrintLine(output.Yellow)

	c.Printer.PrintDetails(fmt.Sprintf("Model: %s, Profile: %s", models.ModelEnum(c.Cache.Model), c.Cache.Profile.Name))
	c.Printer.PrintDetails(fmt.Sprintf("Answered in: %0.2f seconds", res.DurationSeconds()))
	c.Printer.Print("", "")
}
