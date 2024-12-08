package anthropic

import (
	"encoding/json"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

var slugMap = map[string]string{
	"sonnet": "claude-3-5-sonnet-20241022",
}

type AnthropicModel struct {
	model  string
	apiKey string
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

func NewProvider(modelSlug string, apiKey string) *AnthropicModel {

	return &AnthropicModel{
		model:  slugMap[modelSlug],
		apiKey: apiKey,
	}
}

func (m *AnthropicModel) GetUrl() string {
	return "https://api.anthropic.com/v1/messages"
}

func (m *AnthropicModel) GetHeaders() map[string]string {
	return map[string]string{
		"x-api-key":         m.apiKey,
		"anthropic-version": "2023-06-01",
		"Content-Type":      "application/json",
	}
}

func (m *AnthropicModel) PreparePayload(chatHistory *conversation.Conversation, systemPrompt string) (string, error) {

	messages := make([]message, 0, chatHistory.Len())
	for _, qa := range chatHistory.QAs {
		messages = append(messages, message{
			Role:    roleUser,
			Content: qa.GetQuestion(),
		})
		if qa.HasAnswer() {
			messages = append(messages, message{
				Role:    roleAssistant,
				Content: qa.GetAnswer(),
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
