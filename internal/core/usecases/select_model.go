package usecases

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
)

type SelectModelService struct {
	conversationRepo conversation.IConversationRepository
}

func NewSelectModelService(conversationRepo conversation.IConversationRepository) *SelectModelService {
	return &SelectModelService{
		conversationRepo: conversationRepo,
	}
}

func (s *SelectModelService) ListAll() map[models.ModelEnum]string {
	return models.ListModels()
}

func (s *SelectModelService) SetAsActive(model models.ModelEnum) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	conv.SetModelTo(string(model))
	s.conversationRepo.SaveAsActive(conv)

	return nil
}
