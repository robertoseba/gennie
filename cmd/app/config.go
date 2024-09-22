package app

import (
	"strconv"

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

func ConfigProfile(profiles *[]profile.Profile) string {
	menu := menu.NewMenu("Select the profile you want to activate:")
	for idx, profile := range *profiles {
		menu.AddItem(profile.Name, strconv.Itoa(idx))
	}
	selection := menu.Display()

	return selection
}
