package usecases

import (
	"errors"
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

func TestCompleteService(t *testing.T) {
	t.Run("completes the conversation with answers from the API", func(t *testing.T) {
		mockDeps := NewMockDeps()
		mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAIMini.Slug()))
		mockDeps.WithProfile(profile.DefaultProfile())
		mockDeps.WithAPIAnswer("it's an Ai assistant")
		service := mockDeps.createService()

		outputChan, err := service.Execute(&InputDTO{
			Question:     "What is gennie?",
			ProfileSlug:  profile.DefaultProfileSlug,
			Model:        models.OpenAI.Slug(),
			IsFollowUp:   false,
			AppendFile:   "",
			IsStreamable: false,
		})

		answer := <-outputChan
		require.NoError(t, err)
		require.Equal(t, "it's an Ai assistant", answer.Data)
		require.NoError(t, answer.Err)
		mockDeps.AssertExpectations(t)
	})

	// t.Run("appends the content of a file to the conversation", func(t *testing.T) {
	// 	filecontents := "Gennie is an AI assistant that helps you with your daily tasks\n"
	// 	require.NoError(t, os.WriteFile("./testdata_temp.txt", []byte(filecontents), 0644))

	// 	mockDeps := NewMockDeps()
	// 	mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.DefaultModel.Slug()))
	// 	mockDeps.WithProfile(profile.DefaultProfile())
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	outputChan, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "./testdata_temp.txt",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, "What is gennie?\n"+filecontents, returnedConv.LastQuestion())

	// 	os.Remove("./testdata_temp.txt")
	// })

	// t.Run("if model not provided, uses the model from the active conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug()))
	// 	mockDeps.WithProfile(profile.DefaultProfile())
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	returnedConv, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, models.OpenAI.Slug(), returnedConv.ModelSlug)
	// })

	// t.Run("if profile not provided, uses the profile from the active conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	mockDeps.WithActiveConversation(conversation.NewConversation("test-profile", models.OpenAI.Slug()))
	// 	mockDeps.WithProfile(&profile.Profile{
	// 		Slug: "test-profile",
	// 		Name: "test profile",
	// 		Data: "test data",
	// 	})
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	returnedConv, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, "test-profile", returnedConv.ProfileSlug)
	// })
	// t.Run("when model is inputed replaces the model in active conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug()))
	// 	mockDeps.WithProfile(&profile.Profile{
	// 		Slug: "test-profile",
	// 		Name: "test profile",
	// 		Data: "test data",
	// 	})
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	returnedConv, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "test-profile",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, "test-profile", returnedConv.ProfileSlug)
	// })

	// t.Run("when profile is inputed replaces the profile in active conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug()))
	// 	mockDeps.WithProfile(&profile.Profile{
	// 		Slug: "test-profile",
	// 		Name: "test profile",
	// 		Data: "test data",
	// 	})
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	returnedConv, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "test-profile",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, "test-profile", returnedConv.ProfileSlug)
	// })

	// t.Run("when input is a not set as follow up question, creates a new conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	oldConversation := conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug())
	// 	oldConversation.NewQuestion("previous question")
	// 	oldConversation.AnswerLastQuestion("previous answer")
	// 	mockDeps.WithActiveConversation(oldConversation)

	// 	mockDeps.WithProfile(profile.DefaultProfile())
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	outputChan, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "",
	// 		Model:       "",
	// 		IsFollowUp:  false,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, "it's an Ai assistant", <-outputChan)
	// })

	// t.Run("when input is a follow up question, appends the question to the conversation", func(t *testing.T) {
	// 	mockDeps := NewMockDeps()
	// 	oldConversation := conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug())
	// 	oldConversation.NewQuestion("previous question")
	// 	oldConversation.AnswerLastQuestion("previous answer")
	// 	mockDeps.WithActiveConversation(oldConversation)

	// 	mockDeps.WithProfile(profile.DefaultProfile())
	// 	mockDeps.WithAPIAnswer("it's an Ai assistant")
	// 	service := mockDeps.createService()

	// 	returnedConv, err := service.Execute(&InputDTO{
	// 		Question:    "What is gennie?",
	// 		ProfileSlug: "",
	// 		Model:       "",
	// 		IsFollowUp:  true,
	// 		AppendFile:  "",
	// 	})

	// 	require.NoError(t, err)
	// 	mockDeps.AssertExpectations(t)
	// 	require.Equal(t, 2, returnedConv.Len())
	// 	require.Equal(t, "What is gennie?", returnedConv.LastQuestion())
	// 	require.Equal(t, "it's an Ai assistant", returnedConv.LastAnswer())
	// 	require.Equal(t, "previous question", returnedConv.QAs[0].Question.Content)
	// 	require.Equal(t, "previous answer", returnedConv.QAs[0].Answer.Content)
	// })

	// //Error Handling
	t.Run("returns an error if cant find profile", func(t *testing.T) {
		mockDeps := NewMockDeps()
		mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug()))
		mockDeps.WithAPIAnswer("it's an Ai assistant")
		mockDeps.mockProfileRepo.On("FindBySlug", "invalid-profile").Return(nil, errors.New("Invalid"))
		service := mockDeps.createService()

		_, err := service.Execute(&InputDTO{
			Question:    "What is gennie?",
			ProfileSlug: "invalid-profile",
			Model:       "",
			IsFollowUp:  true,
			AppendFile:  "",
		})

		require.ErrorContains(t, err, "Invalid")
	})

	t.Run("returns an error if cant find model", func(t *testing.T) {
		mockDeps := NewMockDeps()
		mockDeps.WithActiveConversation(conversation.NewConversation(profile.DefaultProfileSlug, models.OpenAI.Slug()))
		mockDeps.WithAPIAnswer("it's an Ai assistant")
		mockDeps.WithProfile(profile.DefaultProfile())
		service := mockDeps.createService()

		_, err := service.Execute(&InputDTO{
			Question:    "What is gennie?",
			ProfileSlug: "",
			Model:       "invalid-model",
			IsFollowUp:  false,
			AppendFile:  "",
		})

		require.ErrorIs(t, err, models.ErrModelNotFound)
	})
}

type MockDeps struct {
	mockConvRepo    *mocks.MockConversationRepository
	mockProfileRepo *mocks.MockProfileRepository
	mockApiClient   *apimock.ClientApiMock
}

func NewMockDeps() *MockDeps {
	mockConvRepo := mocks.NewMockConversationRepository()
	mockProfileRepo := mocks.NewMockProfileRepository()
	mockApiClient := apimock.NewClientApiMock()
	mockConvRepo.On("SaveAsActive", mock.Anything).Return(nil)

	return &MockDeps{
		mockConvRepo:    mockConvRepo,
		mockProfileRepo: mockProfileRepo,
		mockApiClient:   mockApiClient,
	}
}

func (m *MockDeps) WithActiveConversation(conv *conversation.Conversation) *MockDeps {
	m.mockConvRepo.On("LoadActive").Return(conv, nil)
	return m
}

func (m *MockDeps) WithProfile(profile *profile.Profile) *MockDeps {
	m.mockProfileRepo.On("FindBySlug", profile.Slug).Return(profile, nil)
	return m
}

func (m *MockDeps) WithAPIAnswer(answer string) *MockDeps {
	mockOpenAIResponse := openaiMock.NewMockOpenAIResponse(answer)
	m.mockApiClient.On("Post", mock.Anything, mock.Anything, mock.Anything).Return(mockOpenAIResponse, nil)
	return m
}

func (m *MockDeps) createService() *CompleteService {
	return NewCompleteService(m.mockConvRepo, m.mockProfileRepo, m.mockApiClient, config.NewConfig())
}

func (m *MockDeps) AssertExpectations(t *testing.T) {
	m.mockConvRepo.AssertExpectations(t)
	m.mockProfileRepo.AssertExpectations(t)
	m.mockApiClient.AssertExpectations(t)
}
