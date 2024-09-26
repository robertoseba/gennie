package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/cmd/app"
	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
	"github.com/robertoseba/gennie/internal/output"
)

func main() {
	cmd.Execute()
	os.Exit(0)

	inputOptions := app.ParseCliOptions()

	c, err := cache.Load()
	if err != nil {
		exitWithError(err)
	}

	if inputOptions.ShowCurrentConfig {
		printer := output.NewPrinter()
		printer.PrintLine(output.Yellow)
		printer.PrintAnswer(fmt.Sprintf("Model: %s ", models.ModelEnum(c.Model)))
		printer.PrintLine(output.Yellow)
		printer.PrintAnswer(fmt.Sprintf("Profile name: %s", c.Profile.Name))
		printer.PrintAnswer(fmt.Sprintf("Profile description: %s", c.Profile.Data))
		printer.PrintLine(output.Yellow)
		os.Exit(0)
	}

	//TODO: change these ifs to a strategy pattern. Ex: Action[inputOptions.Action].Execute()
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

	if c.Profile == nil {
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

	printer := output.NewPrinter()
	printer.PrintLine(output.Yellow)
	printer.PrintAnswer(res.Answer())
	printer.PrintLine(output.Yellow)

	printer.PrintDetails(fmt.Sprintf("Model: %s, Profile: %s", models.ModelEnum(c.Model), c.Profile.Name))
	printer.PrintDetails(fmt.Sprintf("Answered in: %0.2f seconds", res.DurationSeconds()))
	printer.Print("", "")
}

func configModel(c *cache.Cache) error {
	model := output.MenuModel()
	c.SetModel(string(model))
	return c.Save()
}

func configProfile(c *cache.Cache) error {
	profiles, err := profile.LoadProfiles()
	if err != nil {
		exitWithError(err)
	}
	profileIdx := output.MenuProfile(&profiles)

	idx, err := strconv.Atoi(profileIdx)

	if err != nil {
		return err
	}

	c.SetProfile(&profiles[idx])

	return c.Save()
}

func exitWithError(err error) {
	fmt.Fprint(os.Stderr, err)
	fmt.Println()
	os.Exit(1)
}
