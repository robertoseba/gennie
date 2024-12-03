package usecases

import "github.com/robertoseba/gennie/internal/core/conversation"

type ExportConversationService struct {
	conversationRepo conversation.IConversationRepository
}

func NewExportConversationService(conversationRepo conversation.IConversationRepository) *ExportConversationService {
	return &ExportConversationService{
		conversationRepo: conversationRepo,
	}
}

func (s *ExportConversationService) Execute(filename string) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	return s.conversationRepo.ExportToFile(conv, filename)
}
