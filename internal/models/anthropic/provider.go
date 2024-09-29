package anthropic

import (
	"encoding/json"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
)

type AnthropicModel struct {
	url     string
	model   string
	client  httpclient.IHttpClient
	apiKey  string
	headers map[string]string
}

const roleUser = "user"
const roleAssistant = "assistant"

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type prompt struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
}

type content struct {
	ContentType string `json:"type"`
	Text        string `json:"text"`
}
type AnthropicResponse struct {
	Content []content `json:"content"`
}

func NewProvider(modelName string) *AnthropicModel {

	return &AnthropicModel{
		model: modelName,
	}
}

func (m *AnthropicModel) GetUrl() string {
	return "https://api.anthropic.com/v1/messages"
}

func (m *AnthropicModel) GetHeaders() map[string]string {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	return map[string]string{
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
		"Content-Type":      "application/json",
	}
}

func (m *AnthropicModel) PreparePayload(chatHistory *chat.ChatHistory, systemPrompt string) (string, error) {

	messages := make([]message, 0, chatHistory.Len())
	for _, chat := range chatHistory.Responses {
		messages = append(messages, message{
			Role:    roleUser,
			Content: chat.GetQuestion(),
		})
		if chat.HasAnswer() {
			messages = append(messages, message{
				Role:    roleAssistant,
				Content: chat.GetAnswer(),
			})
		}
	}

	p := prompt{
		Model:     m.model,
		Messages:  messages,
		MaxTokens: 1024,
		System:    systemPrompt,
	}
	jsonData, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (m *AnthropicModel) ParseResponse(rawRes []byte) (string, error) {
	var response AnthropicResponse
	err := json.Unmarshal([]byte(rawRes), &response)
	if err != nil {
		return "", err
	}

	return response.Content[0].Text, nil
}

func (m *AnthropicModel) Model() string {
	return m.model
}
