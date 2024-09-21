package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
)

func main() {
	client := httpclient.NewClient()

	inputOptions := app.ParseCliOptions()

	model := models.NewModel(inputOptions.Model, client)

	res, err := model.Ask(inputOptions.Question, nil)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Println()
		os.Exit(1)
		return
	}

	fmt.Printf(res.Answer.Content)

	os.Exit(0)
}
