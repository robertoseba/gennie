package openai

import (
	"os"
	"github.com/robertoseba/gennie/internal/chat"
)

type OpenAIModel struct {
	baseUrl string
	model string
	apiKey string
	promptTemplate string
}

func NewModel() *OpenAIModel {
	return &OpenAIModel{
		baseUrl: "https://api.openai.com/v1/chat/completions",
		model: "gpt4o",
		apiKey: os.Getenv("OPENAI_API_KEY"),
		promptTemplate: "You: %s\nAI: %s",
	}
}

func (m *OpenAIModel) Ask(question string, history *chat.ChatHistory) chat.Response{
	return chat.Response{}	
}

func (m *OpenAIModel) sendQuestion(test string) string {
	return ""
}