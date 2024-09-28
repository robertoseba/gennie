package cache

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/chat"
)

func Load() (*Cache, error) {
	const cacheFile = ".cache"
	basePath, err := getCacheFolderPath()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(basePath, cacheFile)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Cache{
			Model:         "",
			Profile:       nil,
			cacheFilePath: filePath,
			ChatHistory:   chat.NewChatHistory(),
		}, nil
	}

	cache, err := readFrom(filePath)
	if err != nil {
		return nil, err
	}

	cache.cacheFilePath = filePath

	return cache, nil
}

func getCacheFolderPath() (string, error) {
	systemCacheFolder, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheFolder := "gennie"

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

func readFrom(filename string) (*Cache, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache Cache
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, err
	}

	return &cache, nil
}
