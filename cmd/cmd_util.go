package cmd

import (
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

func NewCmdUtil(storage common.IStorage, version string) *CmdUtil {
	return &CmdUtil{
		HttpClient: newHttpClient,
		Printer:    output.NewPrinter(nil, nil),
		Storage:    storage,
		Version:    version,
	}
}

// Client gets instanciated only when needed.
func newHttpClient() httpclient.IHttpClient {
	return httpclient.NewClient()
}
