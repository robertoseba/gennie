package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
)

//go:embed version.txt
var version string

func main() {
	cachePath, err := cache.GetCacheFilePath()
	if err != nil {
		cmd.ExitWithError(err)
	}

	persistence, err := cache.RestoreFrom(cachePath)

	if errors.Is(err, cache.ErrNoCacheFile) {
		persistence = cache.NewCache(cachePath)
		//run cmd config
		fmt.Println("No cache found, let's start config.")
		os.Exit(0)
	}

	if err != nil {
		cmd.ExitWithError(err)
	}

	defer persistence.Save()

	httpClient := httpclient.NewClient()

	printer := output.NewPrinter(nil, nil)

	cmd.NewRootCmd(version, persistence, printer, httpClient).Execute()

}
