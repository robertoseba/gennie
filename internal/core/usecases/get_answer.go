package usecases

import (
	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type GetAnswerService struct {
	conversationRepo conversation.ConversationRepository
	profileRepo      profile.ProfileRepositoryInterface
	apiClient        models.IApiClient
	config           config.Config
}

func NewGetAnswerService(
	cr conversation.ConversationRepository,
	pr profile.ProfileRepositoryInterface,
	apiClient models.IApiClient,
	config config.Config) *GetAnswerService {
	return &GetAnswerService{
		conversationRepo: cr,
		profileRepo:      pr,
		apiClient:        apiClient,
		config:           config,
	}
}

func (s *GetAnswerService) Execute(question string, profileSlugInput string, modelInput string, isFollowUp bool) (*conversation.Conversation, error) {
	var conv *conversation.Conversation
	var err error

	conv, err = s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}

	if profileSlugInput == "" {
		profileSlugInput = conv.ProfileSlug
	}

	if modelInput == "" {
		modelInput = conv.ModelSlug
	}

	if !isFollowUp {
		conv = conversation.NewConversation(profileSlugInput, modelInput)
	}

	conv.NewQuestion(question)
	err = s.completeConversation(conv, profileSlugInput, modelInput)
	if err != nil {
		return nil, err
	}

	err = s.conversationRepo.SaveAsActive(conv)
	if err != nil {
		return conv, err
	}

	return conv, nil
}

func (s *GetAnswerService) completeConversation(conversation *conversation.Conversation, profileSlug, modelSlug string) error {
	profile, err := s.profileRepo.FindBySlug(profileSlug)
	if err != nil {
		return err
	}

	if profileSlug != conversation.ProfileSlug {
		conversation.SetProfileTo(profileSlug)
	}

	model, err := models.NewModel(modelSlug, s.apiClient, s.config)
	if err != nil {
		return err
	}

	if conversation.ModelSlug != modelSlug {
		conversation.SetModelTo(modelSlug)
	}

	err = model.Complete(conversation, profile.Data)

	if err != nil {
		return err
	}

	return nil
}
