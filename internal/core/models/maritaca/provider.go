package maritaca

import (
	"encoding/json"
	"fmt"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

var slugMap = map[string]string{
	"maritaca": "sabia-3",
}

type MaritacaModel struct {
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

func NewProvider(modelSlug string, apiKey string) *MaritacaModel {
	return &MaritacaModel{
		model:  slugMap[modelSlug],
		apiKey: apiKey,
	}
}

func (m *MaritacaModel) GetHeaders() map[string]string {

	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", m.apiKey),
		"Content-Type":  "application/json",
	}
}

func (m *MaritacaModel) GetUrl() string {
	return "https://chat.maritaca.ai/api/chat/completions"
}

func (m *MaritacaModel) PreparePayload(chatHistory *conversation.Conversation, systemPrompt string) (string, error) {
	p := prompt{
		Model: m.model,
		Messages: []message{
			{
				Role:    roleSystem,
				Content: systemPrompt,
			},
		},
	}

	for _, qa := range chatHistory.QAs {
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

func (m *MaritacaModel) ParseResponse(rawRes []byte) (string, error) {
	var response openAiResponse
	err := json.Unmarshal(rawRes, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
