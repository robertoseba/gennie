package app

import (
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
	"github.com/robertoseba/gennie/internal/output/menu"
)

func ConfigModel() models.ModelEnum {
	menu := menu.NewMenu("Select a model:")
	for _, model := range models.ListModels() {
		menu.AddItem(model.String(), string(model))
	}
	selection := menu.Display()

	model := models.ModelEnum(selection)
	return model
}

// returns the slug/filename for the profile
func ConfigProfile() string {
	menu := menu.NewMenu("Select the profile you want to activate:")
	for _, slug := range []string{"default", "linux", "programming"} {
		menu.AddItem(slug, slug)
	}
	selection := menu.Display()

	return selection
}

type Config struct {
	Model   models.ModelEnum
	Profile profile.Profile
}

// TODO: pass persistence as parameter
func LoadConfig() *Config {
	profile := profile.Profile{
		Name:        "Default",
		Description: "Generic default profile",
		Slug:        "default",
		Data:        "You are a cli assistant. You're expert in Linux and programming. You're answer always concise and to the point. If the question is unclear you ask for more information.",
	}

	return &Config{
		Model:   models.OpenAIMini,
		Profile: profile,
	}
}
