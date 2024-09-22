package cache

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/models/profile"
)

const CacheFile = ".cache"

type Cache struct {
	Model   string          `json:"model"`
	Profile profile.Profile `json:"profile"`
}

func (c *Cache) SetModel(modelName string) {
	c.Model = modelName
}

func (c *Cache) SetProfile(profile profile.Profile) {
	c.Profile = profile
}

func LoadCache(basePath string) (*Cache, error) {
	filePath := path.Join(basePath, CacheFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Cache{}, nil
	}

	content, err := readFrom(filePath)
	if err != nil {
		panic(err)
	}

	cache := &Cache{}

	err = json.Unmarshal(content, cache)

	if err != nil {
		panic(err)
	}

	return cache, nil
}

func SaveCache(basePath string, cache *Cache) error {
	filePath := path.Join(basePath, CacheFile)
	err := writeTo(filePath, cache)

	if err != nil {
		panic(err)
	}

	return nil
}

func readFrom(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func writeTo(filename string, cache *Cache) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	_, err = file.Write(content)

	if err != nil {
		return err
	}

	return nil
}
