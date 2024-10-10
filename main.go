package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
)

//go:embed version.txt
var version string

func main() {
	cmdUtil, err := cmd.NewCmdUtil(version)
	if err != nil {
		cmd.ExitWithError(err)
	}
	defer cmdUtil.Storage.Save()

	command := cmd.NewRootCmd(cmdUtil)
	if cmdUtil.Storage.IsNew() {
		command.SetArgs([]string{"config"})
	}

	command.Execute()
}
