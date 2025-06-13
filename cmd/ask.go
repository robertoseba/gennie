package cmd

import (
	"errors"
	"strings"

	"github.com/robertoseba/gennie/internal/core/usecases/complete"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/ui"
	"github.com/spf13/cobra"
)

func NewAskCmd(askCmd *complete.CompleteService, p *output.Printer) *cobra.Command {
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
			isTerminalFlag, _ := cmd.Flags().GetBool("terminal")

			// if isTerminalFlag {
			// 	spinner = output.NewSpinner("Starting...")
			// 	spinner.Start()
			// }
			// startProcessingTime := time.Now()

			req := complete.Request{
				Question:       strings.Join(args, " "),
				ProfileSlug:    profileFlag,
				ModelSlug:      modelFlag,
				IsFollowUp:     isFollowUpFlag,
				AppendFilename: appendFileFlag,
			}

			responseChan, err := askCmd.Execute(cmd.Context(), req)
			if err != nil {
				return err
			}

			if !isTerminalFlag {
				for data := range responseChan {
					if data.Err != nil {
						return data.Err
					}

					if !isTerminalFlag {
						if data.IsApprovalRequest() {
							return errors.New("tool use approval is required. Please run the command in a terminal to approve the request or set in the profile the `requires_approval` flag to false.")
						}
						// When piping we only print the models answer
						if data.IsAnswer() {
							cmd.Print(data.Data)
						}
						continue
					}
				}
				return nil
			}

			return ui.Run(responseChan)

		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

	return cmdAsk
}
