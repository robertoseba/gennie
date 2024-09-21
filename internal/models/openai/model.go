package openai

import (
	"encoding/json"
	"os"

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

func NewModel(client *httpclient.HttpClient, modelName string) *OpenAIModel {
	client.SetBearerToken(os.Getenv("OPEN_API_KEY"))
	return &OpenAIModel{
		url:    "https://api.openai.com/v1/chat/completions",
		model:  modelName,
		client: client,
	}
}

func (m *OpenAIModel) Ask(question string, history *chat.ChatHistory) (*chat.Chat, error) {
	preparedQuestion, err := m.prepareQuestion(question)
	if err != nil {
		return nil, err
	}

	finalResponse := chat.Chat{}
	finalResponse.AddQuestion(question)

	postRes, err := m.client.Post(m.url, preparedQuestion)

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

func (m *OpenAIModel) prepareQuestion(question string) (string, error) {

	//TODO: load preset here for system messages
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
