package usecases

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type SelectProfileService struct {
	profileRepo      profile.ProfileRepository
	conversationRepo conversation.ConversationRepository
}

func NewSelectProfileService(profileRepo profile.ProfileRepository, conversationRepo conversation.ConversationRepository) *SelectProfileService {
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
	return s.conversationRepo.SaveAsActive(conv)
}
