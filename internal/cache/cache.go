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

type Cache struct {
	Config             common.Config
	CachedProfilesPath map[string]string // map[profileSlug]filepath
	ChatHistory        chat.ChatHistory  //Not using pointers so we garantee immutability of chat in cache
	filePath           string
}

func NewCache(filePath string) *Cache {
	return &Cache{
		Config:             common.NewConfig(),
		CachedProfilesPath: map[string]string{},
		ChatHistory:        chat.NewChatHistory(),
		filePath:           filePath,
	}
}

func (c *Cache) Save() error {
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

func (c *Cache) GetConfig() common.Config {
	return c.Config
}

func (c *Cache) SetConfig(config common.Config) {
	c.Config = config
}

func (c *Cache) GetChatHistory() chat.ChatHistory {
	return c.ChatHistory
}

func (c *Cache) SetChatHistory(chatHistory chat.ChatHistory) {
	c.ChatHistory = chatHistory
}

func (c *Cache) GetProfile(profileSlug string) (*profile.Profile, error) {
	if profileSlug == "default" {
		return profile.DefaultProfile(), nil
	}

	filename, ok := c.CachedProfilesPath[profileSlug]
	if !ok {
		return nil, ErrNoProfileSlug
	}

	data, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	return jsonToProfile(data)
}

func (c *Cache) GetProfileSlugs() []string {
	slugs := make([]string, 0, len(c.CachedProfilesPath))
	for k := range c.CachedProfilesPath {
		slugs = append(slugs, k)
	}
	return slugs
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
