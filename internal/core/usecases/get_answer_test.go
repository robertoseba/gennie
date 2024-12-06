package usecases

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	openaiMock "github.com/robertoseba/gennie/internal/core/models/openai/mocks"
	"github.com/robertoseba/gennie/internal/core/profile"
	apimock "github.com/robertoseba/gennie/internal/infra/apiclient/mocks"
	"github.com/robertoseba/gennie/internal/infra/repositories/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAnswerService(t *testing.T) {
	//TODO: create mock IAPIclient
	t.Run("completes the conversation with answers from the API", func(t *testing.T) {
		mockConvRepo := mocks.NewMockConversationRepository()
		conv := conversation.NewConversation("default", models.DefaultModel.Slug())
		mockConvRepo.On("LoadActive").Return(conv, nil)

		expectedConv := conversation.NewConversation("default", models.DefaultModel.Slug())
		expectedConv.NewQuestion("What is gennie?")
		expectedConv.AnswerLastQuestion("it's an Ai assistant")
		mockConvRepo.On("SaveAsActive", expectedConv.LastQuestion(), expectedConv.LastAnswer()).Return(nil)

		mockProfileRepo := mocks.NewMockProfileRepository()
		mockProfileRepo.On("FindBySlug", profile.DefaultProfileSlug).Return(profile.DefaultProfile(), nil)

		config := config.NewConfig()

		mockApiClient := apimock.NewClientApiMock()
		mockOpenAIResponse := openaiMock.NewMockOpenAIResponse("it's an Ai assistant")
		mockApiClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(mockOpenAIResponse, nil)

		service := NewGetAnswerService(mockConvRepo, mockProfileRepo, mockApiClient, *config)
		returnedConv, err := service.Execute(&InputDTO{
			Question:    "What is gennie?",
			ProfileSlug: "",
			Model:       "",
			IsFollowUp:  false,
			AppendFile:  "",
		})

		mockConvRepo.AssertExpectations(t)
		mockApiClient.AssertExpectations(t)
		mockProfileRepo.AssertExpectations(t)

		require.NoError(t, err)
		require.Equal(t, models.DefaultModel.Slug(), returnedConv.ModelSlug)
		require.Equal(t, profile.DefaultProfileSlug, returnedConv.ProfileSlug)
		require.Equal(t, 1, returnedConv.Len())
		require.Equal(t, "What is gennie?", returnedConv.LastQuestion())
		require.Equal(t, "it's an Ai assistant", returnedConv.LastAnswer())
	})

	t.Run("appends the content of a file to the conversation", func(t *testing.T) {
	})

	t.Run("if model not provided, uses the model from the active conversation", func(t *testing.T) {
	})

	t.Run("if profile not provided, uses the profile from the active conversation", func(t *testing.T) {
	})
	t.Run("when model is inputed replaces the model in active conversation", func(t *testing.T) {
	})

	t.Run("when profile is inputed replaces the profile in active conversation", func(t *testing.T) {
	})

	t.Run("when input is a not set as follow up question, creates a new conversation", func(t *testing.T) {
	})

	t.Run("when input is a follow up question, appends the question to the conversation", func(t *testing.T) {
	})

	//Error Handling
	t.Run("returns an error if cant find profile", func(t *testing.T) {
	})

	t.Run("returns an error if cant find model", func(t *testing.T) {
	})

	t.Run("returns an error if cant load active conversation", func(t *testing.T) {
	})
}
