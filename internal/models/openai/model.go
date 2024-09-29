package openai

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
)

type OpenAIModel struct {
	model string
}

const roleUser = "user"
const roleSystem = "system"
const roleAssistant = "assistant"

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type prompt struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type choice struct {
	Message message `json:"message"`
}
type openAiResponse struct {
	Choices []choice `json:"choices"`
}

func NewProvider(modelName string) *OpenAIModel {
	return &OpenAIModel{
		model: modelName,
	}
}

func (m *OpenAIModel) GetHeaders() map[string]string {
	apiKey := os.Getenv("OPEN_API_KEY")

	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
		"Content-Type":  "application/json",
	}
}

func (m *OpenAIModel) GetUrl() string {
	return "https://api.openai.com/v1/chat/completions"
}

func (m *OpenAIModel) PreparePayload(chatHistory *chat.ChatHistory, systemPrompt string) (string, error) {
	p := prompt{
		Model: m.model,
		Messages: []message{
			{
				Role:    roleSystem,
				Content: systemPrompt,
			},
		},
	}

	for _, chat := range chatHistory.Responses {
		p.Messages = append(p.Messages, message{
			Role:    roleUser,
			Content: chat.GetQuestion(),
		})
		if chat.HasAnswer() {
			p.Messages = append(p.Messages, message{
				Role:    roleAssistant,
				Content: chat.GetAnswer(),
			})
		}
	}

	jsonData, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (m *OpenAIModel) ParseResponse(rawRes []byte) (string, error) {
	var response openAiResponse
	err := json.Unmarshal(rawRes, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
