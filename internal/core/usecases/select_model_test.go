package usecases

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
	"github.com/robertoseba/gennie/internal/infra/repositories/mocks"
	"github.com/stretchr/testify/require"
)

func TestModelListAll(t *testing.T) {
	service := NewSelectModelService(nil)
	m := service.ListAll()

	require.Len(t, m, 6)
	require.Equal(t, "GPT-4o-mini (OPENAI)", m[base.OpenAIMini])
}

func TestSetAsActive(t *testing.T) {
	repo := &mocks.MockConversationRepository{}

	activeConv := conversation.NewConversation("profile-slug", "model22")
	activeConv.NewQuestion("question")
	activeConv.AnswerLastQuestion("answer")

	repo.On("LoadActive").Return(activeConv, nil)
	repo.On("SaveAsActive", activeConv).Return(nil)

	service := NewSelectModelService(repo)
	err := service.SetAsActive(base.Groq)

	require.NoError(t, err)
	require.Equal(t, base.Groq.Slug(), activeConv.ModelSlug)
	repo.AssertExpectations(t)
}
