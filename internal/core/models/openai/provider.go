package openai

import (
	"encoding/json"
	"fmt"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

type OpenAIModel struct {
	model  string
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

func NewProvider(modelName string, apiKey string) *OpenAIModel {
	return &OpenAIModel{
		model:  modelName,
		apiKey: apiKey,
	}
}

func (m *OpenAIModel) GetHeaders() map[string]string {

	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", m.apiKey),
		"Content-Type":  "application/json",
	}
}

func (m *OpenAIModel) GetUrl() string {
	return "https://api.openai.com/v1/chat/completions"
}

func (m *OpenAIModel) PreparePayload(conversation *conversation.Conversation, systemPrompt string, isStreamable bool) (string, error) {
	p := prompt{
		Model: m.model,
		Messages: []message{
			{
				Role:    roleSystem,
				Content: systemPrompt,
			},
		},
	}

	for _, qa := range conversation.QAs {
		p.Messages = append(p.Messages, message{
			Role:    roleUser,
			Content: qa.GetQuestion(),
		})
		if qa.HasAnswer() {
			p.Messages = append(p.Messages, message{
				Role:    roleAssistant,
				Content: qa.GetAnswer(),
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

func (m *OpenAIModel) CanStream() bool {
	return false
}

func (m *OpenAIModel) GetStreamParser() func(b []byte) (string, error) {
	return func(b []byte) (string, error) {
		// 	if bytes.Contains(b, []byte("content_block_delta")) && bytes.HasPrefix(b, []byte("data:")) {
		// 		// removes data prefix
		// 		b = bytes.TrimPrefix(b, []byte("data:"))
		// 		var responseData StreamResponse
		// 		err := json.Unmarshal(b, &responseData)

		// 		if err != nil {
		// 			return "", err
		// 		}
		// 		return responseData.Delta.Text, nil
		// 	}
		// 	return "", nil
		return "", nil
	}
}
