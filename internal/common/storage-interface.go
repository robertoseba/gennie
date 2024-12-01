package common

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/profile"
)

type IStorage interface {
	GetConfig() Config
	SetConfig(Config)
	LoadProfileData(string) (*profile.Profile, error)
	GetCachedProfiles() map[string]profile.ProfileInfo
	GetChatHistory() conversation.Conversation
	SetChatHistory(conversation.Conversation)
	GetStorageFilepath() string
	GetCurrProfile() profile.Profile
	SetCurrProfile(profile.Profile)
	GetCurrModelSlug() string
	SetCurrModelSlug(string)
	SetCachedProfiles(map[string]profile.ProfileInfo)
	Clear()
	IsNew() bool
	Save() error
}
