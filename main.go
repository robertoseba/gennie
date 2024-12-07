package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/infra/container"
	"github.com/robertoseba/gennie/internal/output"
)

//go:embed version.txt
var version string

func main() {
	container := container.NewContainer()
	printer := output.NewPrinter(nil, nil)

	command := cmd.NewRootCmd(version, printer, container)
	if container.GetConfig().IsNew() {
		command.SetArgs([]string{"config"})
	}

	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
