package cache

import (
	"encoding/gob"
	"os"
	"path"
)

const cacheFile = ".cache"

func RestoreFrom(filePath string) (*Cache, error) {
	// If file does not exist, returns a new cache
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrNoCacheFile
	}

	cache, err := readFrom(filePath)
	if err != nil {
		return nil, err
	}

	cache.filePath = filePath

	return cache, nil
}

/**
 * It will try to use the system cache directory,
 * if it fails it will fallback to the current directory.
 */
func GetCacheFilePath() (string, error) {

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

	var persistence Cache
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&persistence); err != nil {
		return nil, err
	}

	return &persistence, nil
}

func fallbackPathAsCurrentDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(currentDir, cacheFile), nil
}
