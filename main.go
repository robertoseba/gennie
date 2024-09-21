package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/persistence"
)

func main() {
	inputOptions := app.ParseCliOptions()

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	persistence := persistence.NewDiskPersistence(currentPath)

	if inputOptions.ConfigModel {
		model := app.ConfigModel()
		persistence.Cache.SetModel(string(model))
		persistence.Save()
		os.Exit(0)
	}

	if inputOptions.ConfigProfile {
		profile := app.ConfigProfile()
		persistence.Cache.SetProfile(string(profile))
		persistence.Save()
		os.Exit(0)
	}

	config := app.LoadConfig()

	client := httpclient.NewClient()

	model := models.NewModel(config.Model, client)

	res, err := model.Ask(inputOptions.Question, &config.Profile, nil)

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
