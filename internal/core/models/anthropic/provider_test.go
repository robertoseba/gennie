package anthropic

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/assert"
)

func TestGetHeaders(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	expectedHeaders := map[string]string{
		"x-api-key":         "API_KEY",
		"anthropic-version": "2023-06-01",
		"Content-Type":      "application/json",
	}

	assert.Equal(t, expectedHeaders, m.GetHeaders())
}

func TestGetUrl(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	assert.Equal(t, "https://api.anthropic.com/v1/messages", m.GetUrl())
}

func TestPreparePayload(t *testing.T) {
	m := NewProvider("sonnet", "API_KEY")

	conversation := conversation.NewConversation("test", "test")
	conversation.NewQuestion("Question")
	conversation.AnswerLastQuestion("Answer")
	payload, err := m.PreparePayload(conversation, "System Prompt")

	assert.Nil(t, err)
	assert.JSONEq(t, `{"model":"claude-3-5-sonnet-20241022","messages":[{"role":"user","content":"Question"},{"role":"assistant","content":"Answer"}], "max_tokens":1024, "system":"System Prompt"}`, payload)
}

func TestParseResponse(t *testing.T) {
	m := NewProvider("sonnet", "API_KEY")

	apiResponse := []byte(`{"content":[{"text":"Answer", "type":"text"}]}`)
	modelAnswer, err := m.ParseResponse(apiResponse)

	assert.Nil(t, err)
	assert.Equal(t, "Answer", modelAnswer)
}
