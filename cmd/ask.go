package cmd

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/conversation"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewAskCmd(storage common.IStorage, p *output.Printer, h httpclient.IHttpClient) *cobra.Command {
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

			input := &inputOptions{
				Question:   strings.Join(args, " "),
				IsFollowUp: isFollowUpFlag,
				AppendFile: appendFileFlag,
			}

			if profileFlag != "" {
				profile, err := storage.LoadProfileData(profileFlag)
				if err != nil {
					ExitWithError(err)
				}
				storage.SetCurrProfile(*profile)
			}

			if modelFlag != "" {
				storage.SetCurrModelSlug(modelFlag)
			}

			askModel(storage, p, input, h)

			return nil
		},
	}

	cmdAsk.Flags().BoolVarP(&isFollowUpFlag, "followup", "f", false, "marks the question as a followup question. The previous question will be sent as context.")
	cmdAsk.Flags().StringVarP(&appendFileFlag, "append", "a", "", "appends the content of a file to the question.")
	cmdAsk.Flags().StringVarP(&modelFlag, "model", "m", "", "specifies the model to use.")
	cmdAsk.Flags().StringVarP(&profileFlag, "profile", "p", "", "specifies the profile to use.")

	return cmdAsk
}

type inputOptions struct {
	Question   string
	AppendFile string
	IsFollowUp bool
}

func askModel(storage common.IStorage, p *output.Printer, input *inputOptions, client httpclient.IHttpClient) {
	model := models.NewModel(models.ModelEnum(storage.GetCurrModelSlug()), client, storage.GetConfig())

	chatHistory := storage.GetChatHistory()

	if !input.IsFollowUp && chatHistory.Len() > 0 {
		chatHistory.Clear()
	}

	if input.AppendFile != "" {
		fileContent, err := readFileContents(input.AppendFile)
		if err != nil {
			ExitWithError(err)
		}
		input.Question = fmt.Sprintf("%s\n%s", input.Question, fileContent)

	}

	chat := conversation.NewQA(input.Question)
	chatHistory.AddQA(*chat)

	spinner := output.NewSpinner("Thinking...")
	spinner.Start()
	err := model.CompleteChat(&chatHistory, storage.GetCurrProfile().Data)
	spinner.Stop()

	if err != nil {
		ExitWithError(err)
	}

	lastChat, _ := chatHistory.LastQA()
	storage.SetChatHistory(chatHistory)

	p.PrintLine(output.Yellow)
	p.PrintWithCodeStyling(lastChat.GetAnswer(), output.Yellow)
	p.PrintLine(output.Yellow)

	p.Print(fmt.Sprintf("Model: %s, Profile: %s", models.ModelEnum(storage.GetCurrModelSlug()), storage.GetCurrProfile().Name), output.Cyan)
	p.Print(fmt.Sprintf("Answered in: %0.2f seconds", lastChat.DurationSeconds()), output.Cyan)
	p.Print("", "")
}

func readFileContents(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("File %s not found", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file %s: %s", filePath, err)
	}

	return string(content), nil
}
