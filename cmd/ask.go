package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewAskCmd(askCmd *usecases.CompleteService, p *output.Printer) *cobra.Command {
	var isFollowUpFlag bool
	var appendFileFlag string
	var modelFlag string
	var profileFlag string
	var isStreamableFlag bool
	var isTerminalFlag bool

	cmdAsk := &cobra.Command{
		Use:   "ask [question for the llm model]",
		Short: "You can ask anything here",
		Long:  `The question that will be sent to the llm. If your question contains special characters, please use quotes.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var spinner *output.Spinner

			if !term.IsTerminal(int(os.Stdout.Fd())) || !isTerminalFlag {
				isTerminalFlag = false
				isStreamableFlag = false
			} else {
				spinner = output.NewSpinner("Thinking...")
				spinner.Start()
			}

			startProcessingTime := time.Now()

			dto := &usecases.InputDTO{
				Question:     strings.Join(args, " "),
				ProfileSlug:  profileFlag,
				Model:        modelFlag,
				IsFollowUp:   isFollowUpFlag,
				AppendFile:   appendFileFlag,
				IsStreamable: isStreamableFlag,
			}

			respChan, err := askCmd.Execute(dto)
			if err != nil {
				if isTerminalFlag {
					spinner.Stop()
				}
				return err
			}

			isSpinnerRunning := true
			for d := range respChan {
				if isTerminalFlag && isSpinnerRunning {
					spinner.Stop()
					isSpinnerRunning = false
					p.PrintLine(output.Yellow)
				}

				if d.Err != nil {
					return d.Err
				}
				cmd.Print(d.Data)
			}
			cmd.Println()

			if isTerminalFlag {
				endProcessingTime := time.Now()
				p.PrintLine(output.Yellow)
				p.Print(fmt.Sprintf("Answered in: %0.2f seconds", endProcessingTime.Sub(startProcessingTime).Seconds()), output.Cyan)
				p.Print("", "")
			}

			// p.PrintWithCodeStyling(conversation.LastAnswer(), output.Yellow)
			// p.Print(fmt.Sprintf("Model: %s, Profile: %s", conversation.ModelSlug, conversation.ProfileSlug), output.Cyan)

			return nil
		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")
	cmdAsk.Flags().BoolVarP(&isStreamableFlag, "stream", "s", true, "controls if response should be streamed")
	cmdAsk.Flags().BoolVarP(&isTerminalFlag, "terminal", "t", true, "controls if output if interactive terminal or if should just output plain data. Can be used for piping")

	return cmdAsk
}
