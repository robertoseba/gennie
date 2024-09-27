package anthropic

import (
	"encoding/json"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/profile"
)

type AnthropicModel struct {
	url     string
	model   string
	client  *httpclient.HttpClient
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

func NewModel(client *httpclient.HttpClient, modelName string) *AnthropicModel {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	headers := map[string]string{
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
		"Content-Type":      "application/json",
	}

	return &AnthropicModel{
		url:     "https://api.anthropic.com/v1/messages",
		model:   modelName,
		client:  client,
		apiKey:  apiKey,
		headers: headers,
	}
}

func (m *AnthropicModel) Ask(question string, profile *profile.Profile, history *chat.ChatHistory) (*chat.Chat, error) {
	preparedQuestion, err := m.prepareQuestion(question, profile)
	if err != nil {
		return nil, err
	}

	finalResponse := chat.Chat{}
	finalResponse.AddQuestion(question)

	postRes, err := m.client.Post(m.url, preparedQuestion, m.headers)

	if err != nil {
		return nil, err
	}

	parsedResponse, err := m.parseResponse(postRes)
	if err != nil {
		return nil, err
	}

	finalResponse.AddAnswer(parsedResponse)

	return &finalResponse, nil
}

func (m *AnthropicModel) prepareQuestion(question string, profile *profile.Profile) (string, error) {

	p := prompt{
		Model: m.model,
		Messages: []message{
			{
				Role:    roleUser,
				Content: question,
			},
		},
		MaxTokens: 1024,
		System:    profile.Data,
	}
	jsonData, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (m *AnthropicModel) parseResponse(rawRes []byte) (string, error) {
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
