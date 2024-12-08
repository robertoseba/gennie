package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewAskCmd(askCmd *usecases.CompleteService, p *output.Printer) *cobra.Command {
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
			startProcessingTime := time.Now()
			spinner := output.NewSpinner("Thinking...")
			spinner.Start()

			dto := &usecases.InputDTO{
				Question:    strings.Join(args, " "),
				ProfileSlug: profileFlag,
				Model:       modelFlag,
				IsFollowUp:  isFollowUpFlag,
				AppendFile:  appendFileFlag,
			}

			conversation, err := askCmd.Execute(dto)
			spinner.Stop()
			endProcessingTime := time.Now()

			if err != nil {
				return err
			}

			p.PrintLine(output.Yellow)
			p.PrintWithCodeStyling(conversation.LastAnswer(), output.Yellow)
			p.PrintLine(output.Yellow)

			p.Print(fmt.Sprintf("Model: %s, Profile: %s", conversation.ModelSlug, conversation.ProfileSlug), output.Cyan)
			p.Print(fmt.Sprintf("Answered in: %0.2f seconds", endProcessingTime.Sub(startProcessingTime).Seconds()), output.Cyan)
			p.Print("", "")

			return nil
		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

	return cmdAsk
}
