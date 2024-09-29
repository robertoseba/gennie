package cache

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/profile"
)

type Cache struct {
	Model            string
	Profile          *profile.Profile
	FilePath         string
	ChatHistory      *chat.ChatHistory
	ProfileFilenames map[string]string // map[profileSlug]filename
}

func (c *Cache) SetModel(modelName string) {
	c.Model = modelName
}

func (c *Cache) SetProfile(profile *profile.Profile) {
	c.Profile = profile
}

func (c *Cache) Save() error {
	if c.FilePath == "" {
		return nil
	}

	file, err := os.Create(c.FilePath)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Clear() error {
	if c.FilePath == "" {
		return nil
	}

	if _, err := os.Stat(c.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("Cache has not been created yet")
	}

	c.resetCache()

	return os.Remove(c.FilePath)
}

func (c *Cache) resetCache() {
	c.Model = ""
	c.Profile = nil
	c.ChatHistory = chat.NewChatHistory()
}
