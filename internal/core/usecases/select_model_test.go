package usecases

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/infra/repositories/mocks"
	"github.com/stretchr/testify/require"
)

func TestModelListAll(t *testing.T) {
	service := NewSelectModelService(nil)
	m := service.ListAll()

	require.Len(t, m, 6)
	require.Equal(t, "GPT-4o-mini (OPENAI)", m[models.OpenAIMini])
}

func TestSetAsActive(t *testing.T) {
	repo := &mocks.MockConversationRepository{}

	activeConv := conversation.NewConversation("profile-slug", "model22")
	repo.On("LoadActive").Return(activeConv, nil)
	repo.On("SaveAsActive", activeConv).Return(nil)

	service := NewSelectModelService(repo)
	err := service.SetAsActive(models.Groq)

	require.NoError(t, err)
	require.Equal(t, models.Groq.Slug(), activeConv.ModelSlug)
	repo.AssertExpectations(t)
}
