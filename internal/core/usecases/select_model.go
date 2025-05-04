package usecases

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
)

type SelectModelService struct {
	conversationRepo conversation.IConversationRepository
}

func NewSelectModelService(conversationRepo conversation.IConversationRepository) *SelectModelService {
	return &SelectModelService{
		conversationRepo: conversationRepo,
	}
}

func (s *SelectModelService) ListAll() map[base.ModelEnum]string {
	return base.ListModels()
}

func (s *SelectModelService) SetAsActive(model base.ModelEnum) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	conv.SetModelTo(string(model))
	return s.conversationRepo.SaveAsActive(conv)
}
