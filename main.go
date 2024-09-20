package main

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/httpclient"
)

func main() {
	client := httpclient.NewClient()
	res, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}
