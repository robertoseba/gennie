package anthropic

import (
	"bytes"
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
	Stream    bool      `json:"stream"`
}

type content struct {
	ContentType string `json:"type"`
	Text        string `json:"text"`
}
type AnthropicResponse struct {
	Content []content `json:"content"`
}

// StreamResponse represents the structure of Anthropic API streaming responses
type deltaResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type streamResponse struct {
	Type  string        `json:"type"`
	Index int           `json:"index,omitempty"`
	Delta deltaResponse `json:"delta,omitempty"`
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

func (m *AnthropicModel) PreparePayload(conv *conversation.Conversation, systemPrompt string, isStreamable bool) (string, error) {

	messages := make([]message, 0, conv.Len())
	for _, qa := range conv.QAs {
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
		Stream:    isStreamable,
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

func (m *AnthropicModel) CanStream() bool {
	return true
}

func (m *AnthropicModel) GetStreamParser() func(b []byte) (string, error) {
	return func(b []byte) (string, error) {
		if bytes.Contains(b, []byte("content_block_delta")) && bytes.HasPrefix(b, []byte("data:")) {
			// removes data prefix
			b = bytes.TrimPrefix(b, []byte("data:"))
			var responseData streamResponse
			err := json.Unmarshal(b, &responseData)

			if err != nil {
				return "", err
			}
			return responseData.Delta.Text, nil
		}
		return "", nil
	}
}
