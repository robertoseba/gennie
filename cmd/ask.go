package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/usecases/complete"
	"github.com/robertoseba/gennie/internal/output"
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
			var spinner *output.Spinner
			isTerminalFlag, _ := cmd.Flags().GetBool("terminal")

			if isTerminalFlag {
				spinner = output.NewSpinner("Starting...")
				spinner.Start()
			}
			startProcessingTime := time.Now()

			req := complete.Request{
				Question:       strings.Join(args, " "),
				ProfileSlug:    profileFlag,
				ModelSlug:      modelFlag,
				IsFollowUp:     isFollowUpFlag,
				AppendFilename: appendFileFlag,
			}

			responseChan, err := askCmd.Execute(cmd.Context(), req)
			if err != nil {
				if isTerminalFlag {
					spinner.Stop()
				}
				return err
			}

			var modelInfo, profileInfo string

			for data := range responseChan {
				if data.Err != nil {
					return data.Err
				}

				if !isTerminalFlag {
					if data.IsApprovalRequest() {
						return errors.New("Tool use approval is required. Please run the command in a terminal to approve the request or set in the profile the `requires_approval` flag to false.")
					}
					// When piping we only print the models answer
					if data.IsAnswer() {
						cmd.Print(data.Data)
					}
					continue
				}

				switch data.Type {
				case complete.RtLoading:
					if !spinner.IsRunning() {
						spinner.Start()
					}
					spinner.SetMessage(data.Data)

				case complete.RtModel:
					modelInfo = data.Data

				case complete.RtProfile:
					profileInfo = data.Data

				case complete.RtApprovalReq:
					if spinner.IsRunning() {
						spinner.Stop()
					}
					cmd.Println(data.Data)
					cmd.Println("Please approve the tool use to continue. Should I continue? (y/N)")
					shouldContinue, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')
					if strings.TrimSpace(string(shouldContinue)) != "y" {
						cmd.Println("Tool use approval was not granted. Exiting.")
						return nil
					}

				default:
					if spinner.IsRunning() {
						spinner.Stop()
					}
					cmd.Print(data.Data)
				}
			}

			if isTerminalFlag {
				endProcessingTime := time.Now()
				cmd.Println()
				p.PrintLine(output.Yellow)
				cmd.Printf("Answered in: %0.2f seconds\n", endProcessingTime.Sub(startProcessingTime).Seconds())
				cmd.Printf("Model: %s | Profile: %s\n\n", modelInfo, profileInfo)
			}

			return nil
		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

	return cmdAsk
}
