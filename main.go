package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
)

//go:embed version.txt
var version string

func main() {
	cache, err := cache.Load()
	if err != nil {
		cmd.ExitWithError(err)
	}

	httpClient := httpclient.NewClient()

	printer := output.NewPrinter(nil, nil)

	cmd.NewRootCmd(version, cache, printer, httpClient).Execute()

}
