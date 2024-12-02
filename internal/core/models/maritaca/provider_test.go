package maritaca

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/assert"
)

func TestGetHeaders(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	expectedHeaders := map[string]string{
		"Authorization": "Bearer API_KEY",
		"Content-Type":  "application/json",
	}

	assert.Equal(t, expectedHeaders, m.GetHeaders())
}

func TestGetUrl(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	assert.Equal(t, "https://conversation.maritaca.ai/api/chat/completions", m.GetUrl())
}

func TestPreparePayload(t *testing.T) {
	m := NewProvider("maritaca", "API_KEY")

	conversation := conversation.NewConversation("test", "test")
	conversation.NewQuestion("Question")
	conversation.AnswerLastQuestion("Answer")
	payload, err := m.PreparePayload(conversation, "System Prompt")

	assert.Nil(t, err)
	assert.JSONEq(t, `{"model":"sabia-3","messages":[{"role":"system","content":"System Prompt"},{"role":"user","content":"Question"},{"role":"assistant","content":"Answer"}]}`, payload)
}

func TestParseResponse(t *testing.T) {
	m := NewProvider("gpt-4o", "API_KEY")

	apiResponse := []byte(`{"choices":[{"message":{"role":"assistant","content":"Answer"}}]}`)
	modelAnswer, err := m.ParseResponse(apiResponse)

	assert.Nil(t, err)
	assert.Equal(t, "Answer", modelAnswer)
}
