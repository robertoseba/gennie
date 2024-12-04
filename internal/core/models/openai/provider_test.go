package openai

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/require"
)

func TestGetHeaders(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	expectedHeaders := map[string]string{
		"Authorization": "Bearer API_KEY",
		"Content-Type":  "application/json",
	}

	require.Equal(t, expectedHeaders, m.GetHeaders())
}

func TestGetUrl(t *testing.T) {
	m := NewProvider("test", "API_KEY")
	require.Equal(t, "https://api.openai.com/v1/chat/completions", m.GetUrl())
}

func TestPreparePayload(t *testing.T) {
	m := NewProvider("gpt-4o", "API_KEY")

	conversation := conversation.NewConversation("test", "test")
	conversation.NewQuestion("Question")
	conversation.AnswerLastQuestion("Answer")
	payload, err := m.PreparePayload(conversation, "System Prompt")

	require.NoError(t, err)
	require.JSONEq(t, `{"model":"gpt-4o","messages":[{"role":"system","content":"System Prompt"},{"role":"user","content":"Question"},{"role":"assistant","content":"Answer"}]}`, payload)
}

func TestParseResponse(t *testing.T) {
	m := NewProvider("gpt-4o", "API_KEY")

	apiResponse := []byte(`{"choices":[{"message":{"role":"assistant","content":"Answer"}}]}`)
	modelAnswer, err := m.ParseResponse(apiResponse)

	require.NoError(t, err)
	require.Equal(t, "Answer", modelAnswer)
}
