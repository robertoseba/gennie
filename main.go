package main

import (
	_ "embed"
	"errors"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
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
		if errors.Is(err, cache.ErrNoStoreFile) {
			storage = cache.NewStorage(storagePath)
			shouldConfig = true
		} else {
			cmd.ExitWithError(err)
		}
	}

	defer storage.Save()

	httpClient := httpclient.NewClient()
	printer := output.NewPrinter(nil, nil)

	command := cmd.NewRootCmd(version, storage, printer, httpClient)
	if shouldConfig {
		command.SetArgs([]string{"config"})
	}
	command.Execute()

}
