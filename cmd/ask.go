package cmd

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/llmcore"
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

			req := &complete.Request{
				Question:       strings.Join(args, " "),
				ProfileSlug:    profileFlag,
				ModelSlug:      modelFlag,
				IsFollowUp:     isFollowUpFlag,
				AppendFilename: appendFileFlag,
			}

			ctx := complete.AddRequestToCtx(cmd.Context(), req)
			respChan, err := askCmd.Execute(ctx)
			if err != nil {
				if isTerminalFlag {
					spinner.Stop()
				}
				return err
			}

			var modelInfo, profileInfo string

			for response := range respChan {
				if response.Err != nil {
					return response.Err
				}

				if !isTerminalFlag {
					// When piping we only print the models answer
					if response.Type == "" || response.Type == llmcore.LlmAnswer {
						cmd.Print(response.Data)
					}
					if response.Type == llmcore.ApprovalRequest {
						return errors.New("Tool use approval is required. Please run the command in a terminal to approve the request.")
					}
					continue
				}

				if spinner.IsRunning() {
					switch response.Type {
					case llmcore.LoadingInfo:
						spinner.SetMessage(response.Data)
					case llmcore.ModelInfo:
						modelInfo = response.Data
					case llmcore.ProfileInfo:
						profileInfo = response.Data
					case llmcore.ApprovalRequest:
						spinner.Stop()
						cmd.Println("Please approve the tool use to continue:")
						cmd.Println(response.Data)
						cmd.Println("If you want to cancel the request, please use Ctrl+C.")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
						spinner.Start()
					default: // Answer received
						spinner.Stop()
						cmd.Print(response.Data)
					}
				} else {
					switch response.Type {
					case llmcore.ApprovalRequest:
						cmd.Println()
						p.Print("\nPlease approve the tool use to continue:", output.Yellow)
						cmd.Println(response.Data)
						p.Print("\nIf you want to cancel the request, please use Ctrl+C.", output.Red)
						bufio.NewReader(os.Stdin).ReadBytes('\n')
					case llmcore.LoadingInfo:
						cmd.Println()
						spinner = output.NewSpinner(response.Data)
						spinner.Start()
					default:
						cmd.Print(response.Data)
					}
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
