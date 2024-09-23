package output

import (
	"strconv"

	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/models/profile"
)

// TODO: create in menu a NewMenu that receivees titles and items[names, values] and returns a menu
func MenuModel() models.ModelEnum {
	menu := NewMenu("Select a model:")
	for _, model := range models.ListModels() {
		menu.AddItem(model.String(), string(model))
	}
	selection := menu.Display()

	model := models.ModelEnum(selection)
	return model
}

func MenuProfile(profiles *[]profile.Profile) string {
	menu := NewMenu("Select the profile you want to activate:")
	for idx, profile := range *profiles {
		menu.AddItem(profile.Name, strconv.Itoa(idx))
	}
	selection := menu.Display()

	return selection
}
