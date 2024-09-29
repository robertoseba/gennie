package cache

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/chat"
)

const cacheFile = ".cache"

func Load() (*Cache, error) {
	filePath, err := getCacheFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &Cache{
			Model:       "",
			Profile:     nil,
			FilePath:    filePath,
			ChatHistory: chat.NewChatHistory(),
		}, nil
	}

	cache, err := readFrom(filePath)
	if err != nil {
		return nil, err
	}

	cache.FilePath = filePath

	return cache, nil
}

/**
 * It will try to use the system cache directory,
 * if it fails it will fallback to the current directory.
 */
func getCacheFilePath() (string, error) {

	//TODO: Add support for env var to set cache dir

	systemCacheDir, err := os.UserCacheDir()
	if err != nil || systemCacheDir == "" {
		return fallbackPathAsCurrentDir()
	}

	cacheDirName := "gennie"
	cacheDirPath := path.Join(systemCacheDir, cacheDirName)

	if _, err := os.Stat(cacheDirPath); os.IsNotExist(err) {

		err = os.Mkdir(cacheDirPath, 0755)
		if err != nil {
			return fallbackPathAsCurrentDir()
		}
	}

	return path.Join(cacheDirPath, cacheFile), nil
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

func fallbackPathAsCurrentDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(currentDir, cacheFile), nil
}
