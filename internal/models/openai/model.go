package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
)

type OpenAIModel struct {
	url     string
	model   string
	client  httpclient.IHttpClient
	apiKey  string
	headers map[string]string
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

func NewModel(client httpclient.IHttpClient, modelName string) *OpenAIModel {
	apiKey := os.Getenv("OPEN_API_KEY")
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
		"Content-Type":  "application/json",
	}

	return &OpenAIModel{
		url:     "https://api.openai.com/v1/chat/completions",
		model:   modelName,
		client:  client,
		apiKey:  apiKey,
		headers: headers,
	}
}

func (m *OpenAIModel) CompleteChat(chatHistory *chat.ChatHistory, systemPrompt string) error {
	lastChat, ok := chatHistory.LastChat()
	if !ok {
		return errors.New("Chat history is empty")
	}

	if lastChat.GetAnswer() != "" {
		return errors.New("Last chat is already completed with answer")
	}

	payload, err := m.preparePayload(chatHistory, systemPrompt)
	if err != nil {
		return err
	}

	postRes, err := m.client.Post(m.url, payload, m.headers)

	if err != nil {
		return err
	}

	parsedResponse, err := m.parseResponse(postRes)
	if err != nil {
		return err
	}

	chatHistory.SetNewAnswerToLastChat(parsedResponse)

	return nil
}

func (m *OpenAIModel) preparePayload(chatHistory *chat.ChatHistory, systemPrompt string) (string, error) {

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

func (m *OpenAIModel) parseResponse(rawRes []byte) (string, error) {
	var response openAiResponse
	err := json.Unmarshal([]byte(rawRes), &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

func (m *OpenAIModel) Model() string {
	return m.model
}
