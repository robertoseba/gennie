package ollama

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/robertoseba/gennie/internal/conversation"
)

type OllamaAIModel struct {
	model string
	host  string
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

type response struct {
	Message message `json:"message"`
}

func NewProvider(modelName string, host string, model string) *OllamaAIModel {
	host = strings.TrimSuffix(host, "/")

	return &OllamaAIModel{
		model: model,
		host:  host,
	}

}

func (m *OllamaAIModel) GetHeaders() map[string]string {

	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (m *OllamaAIModel) GetUrl() string {
	return fmt.Sprintf("%s/api/chat", m.host)
}

func (m *OllamaAIModel) PreparePayload(chatHistory *conversation.Conversation, systemPrompt string) (string, error) {
	p := prompt{
		Model: m.model,
		Messages: []message{
			{
				Role:    roleSystem,
				Content: systemPrompt,
			},
		},
		Stream: false,
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

func (m *OllamaAIModel) ParseResponse(rawRes []byte) (string, error) {
	var response response
	err := json.Unmarshal(rawRes, &response)
	if err != nil {
		return "", err
	}

	return response.Message.Content, nil
}
