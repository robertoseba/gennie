package openai

import (
	"bytes"
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
	Stream   bool      `json:"stream"`
}

type choice struct {
	Message message `json:"message"`
}
type openAiResponse struct {
	Choices []choice `json:"choices"`
}

type deltaStream struct {
	Delta struct {
		Content string `json:"content,omitempty"`
	} `json:"delta"`
}

type streamResponse struct {
	Choices []deltaStream `json:"choices"`
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
		Stream: isStreamable,
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
	return true
}

func (m *OpenAIModel) GetStreamParser() func(b []byte) (string, error) {
	return func(b []byte) (string, error) {
		if !bytes.HasPrefix(b, []byte("data:")) {
			return "", nil
		}

		if bytes.HasSuffix(b, []byte("[DONE]")) {
			return "", nil
		}

		if !bytes.Contains(b, []byte("choices")) {
			return "", nil
		}

		jsonData := bytes.TrimPrefix(b, []byte("data:"))
		var data streamResponse
		err := json.Unmarshal(jsonData, &data)
		if err != nil {
			return "", err
		}
		return data.Choices[0].Delta.Content, nil
	}
}
