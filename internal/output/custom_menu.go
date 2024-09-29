package output

import (
	"github.com/robertoseba/gennie/internal/models"
	"github.com/robertoseba/gennie/internal/profile"
)

// TODO: create in menu a NewMenu that receivees titles and items[names, values] and returns a menu
func MenuModel(m []models.ModelEnum, selected models.ModelEnum) models.ModelEnum {
	menu := NewMenu("Select a model:")

	idxSelected := 0
	idx := 0
	for _, model := range m {
		if model == selected {
			idxSelected = idx
		}
		menu.AddItem(model.String(), string(model))
		idx++
	}

	selection := menu.Display(idxSelected)

	model := models.ModelEnum(selection)
	return model
}

func MenuProfile(profiles map[string]*profile.Profile, selected string) string {
	menu := NewMenu("Select the profile you want to activate:")
	idxSelected := 0
	idx := 0
	for slug := range profiles {
		if slug == selected {
			idxSelected = idx
		}
		menu.AddItem(profiles[slug].Name, slug)
		idx++
	}
	selection := menu.Display(idxSelected)

	return selection
}
