package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
)

//go:embed version.txt
var version string

func main() {
	//TODO: loadConfig(configDir)
	//TODO: Build dependencies and inject them (repositories and httpclient)

	//TODO: check if config is new?

	command := cmd.NewRootCmd(cmdUtil)
	if cmdUtil.Storage.IsNew() {
		command.SetArgs([]string{"config"})
	}

	command.Execute()
}
