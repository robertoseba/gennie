package openai

import (
	"encoding/json"
	"os"
	"time"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
)

type OpenAIModel struct {
	url    string
	model  string
	client *httpclient.HttpClient
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

func NewModel(client *httpclient.HttpClient) *OpenAIModel {
	client.SetBearerToken(os.Getenv("OPEN_API_KEY"))
	return &OpenAIModel{
		url:    "https://api.openai.com/v1/chat/completions",
		model:  "gpt4o",
		client: client,
	}
}

func (m *OpenAIModel) Ask(question string, history *chat.ChatHistory) (*chat.Response, error) {
	preparedQuestion, err := m.prepareQuestion(question)
	if err != nil {
		return nil, err
	}

	postRes, err := m.client.Post(m.url, preparedQuestion)

	if err != nil {
		return nil, err
	}

	parsedResponse, err := m.parseResponse(postRes)
	if err != nil {
		return nil, err
	}

	return &chat.Response{
		Question: chat.Input{
			Role:      roleUser,
			Content:   question,
			Timestamp: time.Now(),
		},
		Answer: chat.Output{
			Role:      roleAssistant,
			Content:   parsedResponse,
			Timestamp: time.Now(),
		},
	}, nil
}

func (m *OpenAIModel) sendQuestion(test string) string {
	return ""
}

func (m *OpenAIModel) prepareQuestion(question string) (string, error) {

	p := prompt{
		Model: "gpt-4o-mini",
		Messages: []message{
			{
				Role:    roleSystem,
				Content: "you are a helpful cli assistant expert in linux and programming. Please answer short and concise answers. ",
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
