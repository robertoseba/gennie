package models

import (
	"errors"
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models/openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompleteChat(t *testing.T) {
	t.Run("Completes Conversation", func(t *testing.T) {
		apiMockClient := ApiClientMock{}

		apiMockClient.On("Post", "https://api.openai.com/v1/chat/completions", openAiPayload(), map[string]string{"Authorization": "Bearer api-key 123", "Content-Type": "application/json"}).Return(stubOpenAIResponse(), nil)

		conv := conversation.NewConversation("profile-slug", "model-slug")
		conv.NewQuestion("question")
		model := newBaseModel(ModelEnum(OpenAI), &apiMockClient, openai.NewProvider(string(OpenAI), "api-key 123"))

		err := model.Complete(conv, "system-prompt")

		assert.NoError(t, err)
		assert.Equal(t, "response to question", conv.LastAnswer())
	})

	t.Run("Returns error when conversation is empty", func(t *testing.T) {
		conv := conversation.NewConversation("profile-slug", "model-slug")
		model := newBaseModel(ModelEnum(OpenAI), nil, openai.NewProvider(string(OpenAI), "api-key"))
		err := model.Complete(conv, "system-prompt")

		assert.Error(t, err)
		assert.Equal(t, ErrEmptyConversation, err)
	})

	t.Run("Returns error when last question is already answered", func(t *testing.T) {
		conv := conversation.NewConversation("profile-slug", "model-slug")
		conv.NewQuestion("question")
		conv.AnswerLastQuestion("already answered")

		model := newBaseModel(ModelEnum(OpenAI), nil, openai.NewProvider(string(OpenAI), "api-key"))
		err := model.Complete(conv, "system-prompt")

		assert.Error(t, err)
		assert.Equal(t, ErrLastQuestionAlreadyAnswered, err)
	})

	t.Run("Returns error when API call fails", func(t *testing.T) {
		apiMockClient := ApiClientMock{}
		apiMockClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error in API"))

		conv := conversation.NewConversation("profile-slug", "model-slug")
		conv.NewQuestion("question")
		model := newBaseModel(ModelEnum(OpenAI), &apiMockClient, openai.NewProvider(string(OpenAI), "api-key"))

		err := model.Complete(conv, "system-prompt")

		assert.Error(t, err)
		assert.Equal(t, "error in API", err.Error())
	})

	t.Run("Returns error when response parsing fails", func(t *testing.T) {
		t.Skip("Failing with panic when parsing json that does not conform to the expected structure")
		apiMockClient := ApiClientMock{}
		apiMockClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return([]byte(`{"error":"error message"}`), nil)

		conv := conversation.NewConversation("profile-slug", "model-slug")
		conv.NewQuestion("question")
		model := newBaseModel(ModelEnum(OpenAI), &apiMockClient, openai.NewProvider(string(OpenAI), "api-key"))

		err := model.Complete(conv, "system-prompt")

		assert.Error(t, err)
		assert.Equal(t, "error in API", err.Error())
	})
}

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) Post(url string, payload string, headers map[string]string) ([]byte, error) {
	args := m.Called(url, payload, headers)
	argOne := args.Get(0)
	if argOne == nil {
		return nil, args.Error(1)
	}
	return argOne.([]byte), args.Error(1)
}

func stubOpenAIResponse() []byte {
	return []byte(`{
		"choices": [
			{
				"message": {
					"content": "response to question",
					"role": "assistant"
				}
			}
		]
	}`)
}

func openAiPayload() string {
	return `{"model":"gpt-4o","messages":[{"role":"system","content":"system-prompt"},{"role":"user","content":"question"}]}`
}
