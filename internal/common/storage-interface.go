package common

import (
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/profile"
)

type IStorage interface {
	GetConfig() Config
	SetConfig(Config)
	LoadProfileData(string) (*profile.Profile, error)
	GetCachedProfiles() map[string]profile.ProfileInfo
	GetChatHistory() chat.ChatHistory
	SetChatHistory(chat.ChatHistory)
	GetStorageFilepath() string
	GetCurrProfile() profile.Profile
	SetCurrProfile(profile.Profile)
	GetCurrModelSlug() string
	SetCurrModelSlug(string)
	Clear()
}
