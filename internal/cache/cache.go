package cache

import (
	"encoding/gob"
	"os"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/models/profile"
)

type Cache struct {
	Model            string
	Profile          *profile.Profile
	cacheFilePath    string
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
	file, err := os.Create(c.cacheFilePath)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(c); err != nil {
		return err
	}
	return nil
}
