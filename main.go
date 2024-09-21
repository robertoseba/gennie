package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
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

	model := models.NewModel(models.OpenAI, client)

	res, err := model.Ask(inputOptions.Question, nil)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Println()
		os.Exit(1)
		return
	}

	fmt.Println(res.Answer.Content)

	os.Exit(0)
}
