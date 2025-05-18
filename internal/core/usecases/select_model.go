package usecases

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmproviders/factory"
)

type SelectModelService struct {
	conversationRepo conversation.IConversationRepository
}

func NewSelectModelService(conversationRepo conversation.IConversationRepository) *SelectModelService {
	return &SelectModelService{
		conversationRepo: conversationRepo,
	}
}

func (s *SelectModelService) ListAll() map[factory.ModelEnum]string {
	return factory.ListModels()
}

func (s *SelectModelService) SetAsActive(model factory.ModelEnum) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	conv.SetModelTo(string(model))
	return s.conversationRepo.SaveAsActive(conv)
}
