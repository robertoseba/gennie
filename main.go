package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
)

//go:embed version.txt
var version string

func main() {
	cmd.Run(version)
}
