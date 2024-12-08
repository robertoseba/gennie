package usecases

import (
	"os"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	apiClient        models.IApiClient
	config           *config.Config
}

type InputDTO struct {
	Question    string
	ProfileSlug string
	Model       string
	IsFollowUp  bool
	AppendFile  string
}

func NewCompleteService(
	cr conversation.IConversationRepository,
	pr profile.IProfileRepository,
	apiClient models.IApiClient,
	config *config.Config) *CompleteService {
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		apiClient:        apiClient,
		config:           config,
	}
}

func (s *CompleteService) Execute(input *InputDTO) (*conversation.Conversation, error) {
	var conv *conversation.Conversation
	var err error

	conv, err = s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}

	if input.ProfileSlug == "" {
		input.ProfileSlug = conv.ProfileSlug
	}

	if input.Model == "" {
		input.Model = conv.ModelSlug
	}

	if !input.IsFollowUp {
		conv = conversation.NewConversation(input.ProfileSlug, input.Model)
	}

	if input.AppendFile != "" {
		content, err := readFile(input.AppendFile)
		if err != nil {
			return nil, err
		}

		input.Question = input.Question + "\n" + content
	}

	if err := conv.NewQuestion(input.Question); err != nil {
		return nil, err
	}

	err = s.completeConversation(conv, input.ProfileSlug, input.Model)
	if err != nil {
		return nil, err
	}

	err = s.conversationRepo.SaveAsActive(conv)
	if err != nil {
		return conv, err
	}

	return conv, nil
}

func (s *CompleteService) completeConversation(conversation *conversation.Conversation, profileSlug, modelSlug string) error {
	profile, err := s.profileRepo.FindBySlug(profileSlug)
	if err != nil {
		return err
	}

	if profileSlug != conversation.ProfileSlug {
		conversation.SetProfileTo(profileSlug)
	}

	model, err := models.NewModel(modelSlug, s.apiClient, *s.config)
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

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}
