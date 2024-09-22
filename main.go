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

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, err := cache.LoadCache(currentPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if c.Model == "" {
		err := configModel(c, currentPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if c.Profile.Slug == "" {
		err := configProfile(c, currentPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if inputOptions.ConfigModel {
		err := configModel(c, currentPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if inputOptions.ConfigProfile {
		err := configProfile(c, currentPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	client := httpclient.NewClient()

	model := models.NewModel(models.ModelEnum(c.Model), client)

	res, err := model.Ask(inputOptions.Question, &c.Profile, nil)

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

func configModel(c *cache.Cache, path string) error {
	model := app.ConfigModel()
	c.SetModel(string(model))
	err := cache.SaveCache(path, c)
	return err
}

func configProfile(c *cache.Cache, path string) error {
	profileSlug := app.ConfigProfile()

	//TODO: load Profile from slug and set it in cache

	c.SetProfile(profile.Profile{
		Name:        "Default",
		Description: "Generic default profile",
		Slug:        profileSlug,
		Data:        "You are a cli assistant. You're expert in Linux and programming. You're answer always concise and to the point. If the question is unclear you ask for more information.",
	})

	err := cache.SaveCache(path, c)
	return err
}
