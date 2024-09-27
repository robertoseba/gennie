package cmd

import (
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/output"
)

type Container struct {
	Cache   *cache.Cache
	Printer *output.Printer
}

func setUp() *Container {

	c, err := cache.Load()
	if err != nil {
		ExitWithError(err)
	}

	return &Container{
		Cache:   c,
		Printer: output.NewPrinter(),
	}
}
