package cache

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/profile"
)

type Storage struct {
	Config         common.Config
	CurrModelSlug  string
	CurrProfile    profile.Profile
	CachedProfiles map[string]profile.ProfileInfo // map[profileSlug]ProfileCache
	ChatHistory    chat.ChatHistory
	filePath       string
}

func NewStorage(filePath string) *Storage {
	return &Storage{
		Config:         common.NewConfig(),
		CurrModelSlug:  "default",
		CurrProfile:    *profile.DefaultProfile(),
		CachedProfiles: map[string]profile.ProfileInfo{},
		ChatHistory:    chat.NewChatHistory(),
		filePath:       filePath,
	}
}

func (c *Storage) Save() error {
	if c.filePath == "" {
		return nil
	}

	file, err := os.Create(c.filePath)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		return err
	}
	return nil
}

func (c *Storage) GetCurrModelSlug() string {
	return c.CurrModelSlug
}

func (c *Storage) SetCurrModelSlug(slug string) {
	c.CurrModelSlug = slug
}

func (c *Storage) GetCurrProfile() profile.Profile {
	return c.CurrProfile
}

func (c *Storage) SetCurrProfile(profile profile.Profile) {
	c.CurrProfile = profile
}

func (c *Storage) GetStorageFilepath() string {
	return c.filePath
}

func (c *Storage) GetConfig() common.Config {
	return c.Config
}

func (c *Storage) SetConfig(config common.Config) {
	c.Config = config
}

func (c *Storage) GetChatHistory() chat.ChatHistory {
	return c.ChatHistory
}

func (c *Storage) SetChatHistory(chatHistory chat.ChatHistory) {
	c.ChatHistory = chatHistory
}

func (c *Storage) SetCachedProfiles(profiles map[string]profile.ProfileInfo) {
	c.CachedProfiles = profiles
}

func (c *Storage) LoadProfileData(profileSlug string) (*profile.Profile, error) {
	if profileSlug == "default" {
		return profile.DefaultProfile(), nil
	}

	profileInfo, ok := c.CachedProfiles[profileSlug]
	if !ok {
		return nil, ErrNoProfileSlug
	}

	data, err := loadFile(profileInfo.Filepath)
	if err != nil {
		return nil, err
	}

	return jsonToProfile(data)
}

func (c *Storage) Clear() {
	c.CurrProfile = *profile.DefaultProfile()
	c.ChatHistory = chat.NewChatHistory()
	c.Config = common.NewConfig()
	c.CurrModelSlug = "default"
	c.CachedProfiles = map[string]profile.ProfileInfo{}
}

func (c *Storage) GetCachedProfiles() map[string]profile.ProfileInfo {
	return c.CachedProfiles
}

func loadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, err
}

func jsonToProfile(data []byte) (*profile.Profile, error) {
	profile := &profile.Profile{}

	err := json.Unmarshal(data, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}
