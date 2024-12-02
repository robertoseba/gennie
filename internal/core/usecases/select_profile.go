package usecases

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type SelectProfileService struct {
	profileRepo      profile.ProfileRepositoryInterface
	conversationRepo conversation.ConversationRepository
}

func NewSelectProfileService(profileRepo profile.ProfileRepositoryInterface, conversationRepo conversation.ConversationRepository) *SelectProfileService {
	return &SelectProfileService{
		profileRepo:      profileRepo,
		conversationRepo: conversationRepo,
	}
}

func (s *SelectProfileService) ListAll() (map[string]*profile.Profile, error) {
	return s.profileRepo.ListAll()
}

func (s *SelectProfileService) SetAsActive(profile *profile.Profile) error {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return err
	}
	conv.SetProfileTo(profile.Slug)
	s.conversationRepo.SaveAsActive(conv)

	return nil
}
