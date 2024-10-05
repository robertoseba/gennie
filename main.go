package main

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
)

//go:embed version.txt
var version string

func main() {
	cachePath, err := cache.GetCacheFilePath()
	if err != nil {
		panic(err)
		// cmd.ExitWithError(err)
	}

	c, err := cache.RestoreFrom(cachePath)

	if err != nil {
		if errors.Is(err, cache.ErrNoCacheFile) {
			fmt.Println("No cache found, creating a new one.")
			c = cache.NewCache(cachePath)
		} else {
			panic(err)
			// cmd.ExitWithError(err)
		}
	} else {
		fmt.Printf("Cache loaded from: %s\n", cachePath)
	}
	defer c.Save()

	fmt.Printf("Cache: %v\n", c.Config.CurrProfile)

	// httpClient := httpclient.NewClient()

	// printer := output.NewPrinter(nil, nil)

	// cmd.NewRootCmd(version, cache, printer, httpClient).Execute()

}
