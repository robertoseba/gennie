package cache

import (
	"encoding/gob"
	"os"
	"path"
)

const storageDefaultFile = ".cache"

func RestoreFrom(filePath string) (*Storage, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrNoStoreFile
	}

	store, err := decodeGob(filePath)
	if err != nil {
		return nil, err
	}

	store.filePath = filePath
	store.isNew = false

	return store, nil
}

/**
 * It will try to use the system cache directory,
 * if it fails it will fallback to the current directory.
 */
func GetStorageFilepath() (string, error) {

	systemCacheDir, err := os.UserCacheDir()
	if err != nil || systemCacheDir == "" {
		return fallbackPathAsCurrentDir()
	}

	storageDirName := "gennie"
	storageDirPath := path.Join(systemCacheDir, storageDirName)

	if _, err := os.Stat(storageDirPath); os.IsNotExist(err) {

		err = os.Mkdir(storageDirPath, 0755)
		if err != nil {
			return fallbackPathAsCurrentDir()
		}
	}

	return path.Join(storageDirPath, storageDefaultFile), nil
}

func decodeGob(filename string) (*Storage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var persistence Storage
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

	return path.Join(currentDir, storageDefaultFile), nil
}
