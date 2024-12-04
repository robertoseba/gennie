package ollama

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/require"
)

func TestGetHeaders(t *testing.T) {
	m := NewProvider("test", "host", "model")
	expectedHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	require.Equal(t, expectedHeaders, m.GetHeaders())
}

func TestGetUrl(t *testing.T) {
	m := NewProvider("test", "host-url", "model")
	require.Equal(t, "host-url/api/chat", m.GetUrl())
}

func TestPreparePayload(t *testing.T) {
	m := NewProvider("test", "host", "model")

	conversation := conversation.NewConversation("test", "test")
	conversation.NewQuestion("Question")
	conversation.AnswerLastQuestion("Answer")
	payload, err := m.PreparePayload(conversation, "System Prompt")

	require.NoError(t, err)
	require.JSONEq(t, `{"model":"model","messages":[{"role":"system","content":"System Prompt"},{"role":"user","content":"Question"},{"role":"assistant","content":"Answer"}],"stream":false}`, payload)
}

func TestParseResponse(t *testing.T) {
	m := NewProvider("test", "host", "model")

	apiResponse := []byte(`{"message":{"role":"assistant","content":"Answer"}}`)
	modelAnswer, err := m.ParseResponse(apiResponse)

	require.NoError(t, err)
	require.Equal(t, "Answer", modelAnswer)
}
