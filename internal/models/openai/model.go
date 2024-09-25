package openai

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/profile"
)

type OpenAIModel struct {
	url    string
	model  string
	client *httpclient.HttpClient
	apiKey string
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

func NewModel(client *httpclient.HttpClient, modelName string) *OpenAIModel {
	return &OpenAIModel{
		url:    "https://api.openai.com/v1/chat/completions",
		model:  modelName,
		client: client,
		apiKey: os.Getenv("OPEN_API_KEY"),
	}
}

func (m *OpenAIModel) Ask(question string, profile *profile.Profile, history *chat.ChatHistory) (*chat.Chat, error) {
	preparedQuestion, err := m.prepareQuestion(question, profile)
	if err != nil {
		return nil, err
	}

	finalResponse := chat.Chat{}
	finalResponse.AddQuestion(question)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", m.apiKey),
		"Content-Type":  "application/json",
	}

	postRes, err := m.client.Post(m.url, preparedQuestion, headers)

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

func (m *OpenAIModel) sendQuestion(test string) string {
	return ""
}

func (m *OpenAIModel) prepareQuestion(question string, profile *profile.Profile) (string, error) {

	p := prompt{
		Model: "gpt-4o-mini",
		Messages: []message{
			{
				Role:    roleSystem,
				Content: profile.Data,
			},
			{
				Role:    roleUser,
				Content: question,
			},
		},
	}
	jsonData, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (m *OpenAIModel) parseResponse(rawRes []byte) (string, error) {
	var response openAiResponse
	err := json.Unmarshal([]byte(rawRes), &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
