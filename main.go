package main

import (
	_ "embed"
	"errors"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/cache"
)

//go:embed version.txt
var version string

func main() {
	storagePath, err := cache.GetStorageFilepath()
	if err != nil {
		cmd.ExitWithError(err)
	}

	storage, err := cache.RestoreFrom(storagePath)

	shouldConfig := false

	if err != nil {
		if !errors.Is(err, cache.ErrNoStoreFile) {
			cmd.ExitWithError(err)
		}
		storage = cache.NewStorage(storagePath)
		shouldConfig = true
	}

	defer storage.Save()

	cmdUtil := cmd.NewCmdUtil(storage, version)
	command := cmd.NewRootCmd(cmdUtil)

	if shouldConfig {
		command.SetArgs([]string{"config"})
	}

	command.Execute()
}
