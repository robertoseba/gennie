package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/output"
)

func main() {
	inputOptions := app.ParseCliOptions()

	if inputOptions.ConfigMode {
		model := app.ConfigModel()
		fmt.Println(string(model))
		//TODO: persist model selection
		os.Exit(0)
	}

	//TODO: read model selection from config
	client := httpclient.NewClient()

	model := models.NewModel(models.OpenAIMini, client)

	res, err := model.Ask(inputOptions.Question, nil)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Println()
		os.Exit(1)
		return
	}

	output.PrintAnswer(res.Answer())
	fmt.Printf("\nAnswered in: %0.2f seconds\n\n", res.DurationSeconds())

	os.Exit(0)
}
