package usecases

import "github.com/robertoseba/gennie/internal/core/conversation"

type ExportConversationService struct {
	conversationRepo conversation.ConversationRepository
}

func (s *ExportConversationService) Execute(filename string) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	s.conversationRepo.ExportToFile(conv, filename)
	return nil
}
