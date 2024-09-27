package cmd

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	fmt.Fprint(os.Stderr, err.Error(), "\n")
	os.Exit(1)
}
