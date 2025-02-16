package usecases

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/infra/repositories/mocks"
	"github.com/stretchr/testify/require"
)

func TestConversationService(t *testing.T) {
	// _, repo := setupActiveConversation(t)
	t.Run("Saves conversation to file", func(t *testing.T) {
		mockedRepo := mocks.NewMockConversationRepository()
		conv := conversation.NewConversation("profile-slug", models.DefaultModel.Slug())
		mockedRepo.On("LoadActive").Return(conv, nil)
		mockedRepo.On("ExportToFile", conv, "./conversation.json").Return(nil)

		service := NewConversationService(mockedRepo)
		err := service.SaveTo("./conversation.json")
		require.NoError(t, err)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Loads conversation from file", func(t *testing.T) {
		mockedRepo := mocks.NewMockConversationRepository()
		conv := conversation.NewConversation("profile-slug", models.DefaultModel.Slug())
		mockedRepo.On("LoadFromFile", "./conversation.json").Return(conv, nil)
		mockedRepo.On("SaveAsActive", conv).Return(nil)

		service := NewConversationService(mockedRepo)
		err := service.LoadFrom("./conversation.json")
		require.NoError(t, err)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Retrieves last conversation", func(t *testing.T) {
		mockedRepo := mocks.NewMockConversationRepository()
		conv := conversation.NewConversation("profile-slug", models.DefaultModel.Slug())
		conv.NewQuestion("What is your name?")
		conv.AnswerLastQuestion("My name is Assistant")
		mockedRepo.On("LoadActive").Return(conv, nil)

		service := NewConversationService(mockedRepo)
		convResponse, err := service.LastConversation()

		require.NoError(t, err)
		mockedRepo.AssertExpectations(t)
		require.Equal(t, conv, convResponse)
	})

	t.Run("If last conversation is empty returns null", func(t *testing.T) {
		mockedRepo := mocks.NewMockConversationRepository()
		conv := conversation.NewConversation("profile-slug", models.DefaultModel.Slug())
		mockedRepo.On("LoadActive").Return(conv, nil)

		service := NewConversationService(mockedRepo)
		convResponse, err := service.LastConversation()

		require.NoError(t, err)
		mockedRepo.AssertExpectations(t)
		require.Nil(t, convResponse)
	})
}
