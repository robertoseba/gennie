package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/internal/httpclient"
)

func main() {
	client := httpclient.NewClient()
	res, err := client.Get("http://localhost:8080")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Println()
		os.Exit(1)
		return
	}

	fmt.Println(string(res))
	os.Exit(0)
}
