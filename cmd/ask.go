package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewAskCmd(c *cache.Cache, p *output.Printer, h httpclient.IHttpClient) *cobra.Command {

	var isFollowUpFlag bool
	var appendFileFlag string
	var modelFlag string
	var profileFlag string

	cmdAsk := &cobra.Command{
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

			input := &InputOptions{
				Question:   strings.Join(args, " "),
				IsFollowUp: isFollowUpFlag,
				AppendFile: appendFileFlag,
				Model:      modelFlag,
				Profile:    profileFlag,
			}

			if input.Profile != "" {
				profiles, err := profile.LoadProfiles()
				if err != nil {
					ExitWithError(err)
				}
				p, ok := profiles[input.Profile]

				if !ok {
					ExitWithError(fmt.Errorf("Profile %s not found. Please use gennie profile list to check available profiles.", input.Profile))
				}

				c.SetProfile(p)
			}

			if input.Model == "" && c.Model == "" {
				c.Model = string(models.DefaultModel)
			}

			if input.Profile == "" && c.Profile == nil {
				c.SetProfile(profile.CreateDefaultProfile())
			}

			askModel(c, p, input, h)

			c.Save()

			return nil
		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

	return cmdAsk
}

type InputOptions struct {
	Question   string
	IsFollowUp bool
	AppendFile string
	Model      string
	Profile    string
}

func askModel(c *cache.Cache, p *output.Printer, input *InputOptions, client httpclient.IHttpClient) {

	var model models.IModel

	if input.Model != "" {
		c.SetModel(input.Model)
	}

	model = models.NewModel(models.ModelEnum(c.Model), client)

	if !input.IsFollowUp && c.ChatHistory != nil {
		c.ChatHistory.Clear()
	}

	chat := chat.NewChat(input.Question)
	c.ChatHistory.AddChat(*chat)

	err := model.CompleteChat(c.ChatHistory, c.Profile.Data)

	if err != nil {
		ExitWithError(err)
	}

	p.PrintLine(output.Yellow)
	p.PrintWithCodeStyling(c.ChatHistory.LastAnswer(), output.Yellow)
	p.PrintLine(output.Yellow)

	p.Print(fmt.Sprintf("Model: %s, Profile: %s", models.ModelEnum(model.Model()), c.Profile.Name), output.Cyan)
	p.Print(fmt.Sprintf("Answered in: %0.2f seconds", chat.DurationSeconds()), output.Cyan)
	p.Print("", "")

}
