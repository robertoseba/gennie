package usecases

import "github.com/robertoseba/gennie/internal/core/conversation"

type ConversationService struct {
	conversationRepo conversation.IConversationRepository
}

func NewConversationService(conversationRepo conversation.IConversationRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
	}
}

func (s *ConversationService) SaveTo(filename string) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	return s.conversationRepo.ExportToFile(conv, filename)
}

func (s *ConversationService) LoadFrom(filename string) error {
	conv, err := s.conversationRepo.LoadFromFile(filename)
	if err != nil {
		return err
	}
	return s.conversationRepo.SaveAsActive(conv)
}
