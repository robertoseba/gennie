package cmd

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
	"github.com/robertoseba/gennie/internal/core/usecases/complete"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewAskCmd(askCmd *complete.CompleteService, p *output.Printer) *cobra.Command {
	var isFollowUpFlag bool
	var appendFileFlag string
	var modelFlag string
	var profileFlag string
	var isStreamableFlag bool

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
			} else {
				isStreamableFlag = false
			}

			startProcessingTime := time.Now()

			dto := &complete.InputDTO{
				Question:    strings.Join(args, " "),
				ProfileSlug: profileFlag,
				Model:       modelFlag,
				IsFollowUp:  isFollowUpFlag,
				AppendFile:  appendFileFlag,
			}

			respChan, err := askCmd.Execute(dto)
			if err != nil {
				if isTerminalFlag {
					spinner.Stop()
				}
				return err
			}

			var modelInfo, profileInfo string

			for d := range respChan {
				if d.Err != nil {
					return d.Err
				}

				if !isTerminalFlag {
					if d.Type == "" { // when piping we only print if it's not related to interface types
						cmd.Print(d.Data)
					}
					continue
				}

				if spinner.IsRunning() {
					switch d.Type {
					case base.LoadingInfo:
						spinner.SetMessage(d.Data)
					case base.ModelInfo:
						modelInfo = d.Data
					case base.ProfileInfo:
						profileInfo = d.Data
					case base.ApprovalRequest:
						spinner.Stop()
						cmd.Println("Please approve the tool use to continue:")
						cmd.Println(d.Data)
						cmd.Println("If you want to cancel the request, please use Ctrl+C.")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
						spinner.Start()
					default:
						spinner.Stop()
						cmd.Print(d.Data)
					}
				} else {
					switch d.Type {
					case base.ApprovalRequest:
						cmd.Println()
						p.Print("\nPlease approve the tool use to continue:", output.Yellow)
						cmd.Println(d.Data)
						p.Print("\nIf you want to cancel the request, please use Ctrl+C.", output.Red)
						bufio.NewReader(os.Stdin).ReadBytes('\n')
					case base.LoadingInfo:
						cmd.Println()
						spinner = output.NewSpinner(d.Data)
						spinner.Start()
					default:
						cmd.Print(d.Data)
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
	cmdAsk.Flags().BoolVarP(&isStreamableFlag, "stream", "s", true, "controls if response should be streamed")

	return cmdAsk
}
