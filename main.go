package main

import (
	_ "embed"
	"os"

	"github.com/robertoseba/gennie/cmd"
)

//go:embed version.txt
var version string

func main() {
	cmd.Run(version, os.Stdout, os.Stderr)
}
