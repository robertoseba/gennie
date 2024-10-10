package cmd

import (
	"errors"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/output"
)

type CmdUtil struct {
	HttpClient func() httpclient.IHttpClient
	Printer    *output.Printer
	Storage    common.IStorage
	Version    string
}

func (c *CmdUtil) NewHttpClient() httpclient.IHttpClient {
	return httpclient.NewClient()
}

func NewCmdUtil(version string) (*CmdUtil, error) {
	storage, err := getStorage()
	if err != nil {
		return nil, err
	}
	return &CmdUtil{
		HttpClient: newHttpClient,
		Printer:    output.NewPrinter(nil, nil),
		Storage:    storage,
		Version:    version,
	}, nil
}

// Client gets instanciated only when needed.
func newHttpClient() httpclient.IHttpClient {
	return httpclient.NewClient()
}

func getStorage() (common.IStorage, error) {
	storagePath, err := cache.GetStorageFilepath()
	if err != nil {
		return nil, err
	}

	storage, err := cache.RestoreFrom(storagePath)
	if err != nil {
		if !errors.Is(err, cache.ErrNoStoreFile) {
			return nil, err
		}
		storage = cache.NewStorage(storagePath)
	}

	return storage, nil
}
