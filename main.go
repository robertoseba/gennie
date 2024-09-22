package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
	"github.com/robertoseba/gennie/internal/output"
)

func main() {
	inputOptions := app.ParseCliOptions()

	c, err := cache.Load()
	if err != nil {
		exitWithError(err)
	}

	if inputOptions.ConfigModel {
		err := configModel(c)
		if err != nil {
			exitWithError(err)
		}
		os.Exit(0)
	}

	if inputOptions.ConfigProfile {
		err := configProfile(c)
		if err != nil {
			exitWithError(err)
		}
		os.Exit(0)
	}

	if c.Model == "" {
		err := configModel(c)
		if err != nil {
			exitWithError(err)
		}
	}

	if c.Profile == nil || c.Profile.Slug == "" {
		err := configProfile(c)
		if err != nil {
			exitWithError(err)
		}
	}

	client := httpclient.NewClient()

	model := models.NewModel(models.ModelEnum(c.Model), client)

	res, err := model.Ask(inputOptions.Question, c.Profile, nil)

	if err != nil {
		exitWithError(err)
	}

	output.PrintAnswer(res.Answer())
	fmt.Printf("\n\n%sModel: %s, Profile: %s%s\n", output.Gray, c.Model, c.Profile.Name, output.Reset)
	fmt.Printf("%sAnswered in: %0.2f seconds%s\n\n", output.Gray, res.DurationSeconds(), output.Reset)
}

func configModel(c *cache.Cache) error {
	model := app.ConfigModel()
	c.SetModel(string(model))
	return c.Save()
}

func configProfile(c *cache.Cache) error {
	profileSlug := app.ConfigProfile()

	//TODO: load Profile from slug and set it in cache
	c.SetProfile(&profile.Profile{
		Name:        "Default",
		Description: "Generic default profile",
		Slug:        profileSlug,
		Data:        "You are a cli assistant. You're expert in Linux and programming. You're answer always concise and to the point. If the question is unclear you ask for more information.",
	})

	return c.Save()
}

func exitWithError(err error) {
	fmt.Fprint(os.Stderr, err)
	fmt.Println()
	os.Exit(1)
}
