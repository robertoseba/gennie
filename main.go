package main

import (
	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
)

func main() {
	cache, err := cache.Load()
	if err != nil {
		cmd.ExitWithError(err)
	}

	httpClient := httpclient.NewClient()

	cmd.NewRootCmd(cache, output.NewPrinter(nil, nil), httpClient).Execute()

}
