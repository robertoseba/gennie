package cache

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/models/profile"
)

type Cache struct {
	Model    string           `json:"model"`
	Profile  *profile.Profile `json:"profile"`
	filePath string
}

func (c *Cache) SetModel(modelName string) {
	c.Model = modelName
}

func (c *Cache) SetProfile(profile *profile.Profile) {
	c.Profile = profile
}

func (c *Cache) Save() error {
	err := writeTo(c.filePath, c)

	if err != nil {
		return err
	}

	return nil
}

func Load() (*Cache, error) {
	const cacheFile = ".cache"
	basePath, err := getCacheFolderPath()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(basePath, cacheFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Cache{
			Model:    "",
			Profile:  nil,
			filePath: filePath,
		}, nil
	}

	content, err := readFrom(filePath)
	if err != nil {
		return nil, err
	}

	cache := &Cache{}

	err = json.Unmarshal(content, cache)

	cache.filePath = filePath

	if err != nil {
		return nil, err
	}

	return cache, nil
}

func getCacheFolderPath() (string, error) {
	systemCacheFolder, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheFolder := "ginnie"

	if _, err := os.Stat(systemCacheFolder); os.IsNotExist(err) {
		curr, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return curr, nil
	}

	if _, err := os.Stat(path.Join(systemCacheFolder, cacheFolder)); os.IsNotExist(err) {

		err = os.Mkdir(path.Join(systemCacheFolder, cacheFolder), 0755)
		if err != nil {
			return "", err
		}
	}

	return path.Join(systemCacheFolder, cacheFolder), nil
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
