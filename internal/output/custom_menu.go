package output

import "github.com/robertoseba/gennie/internal/core/models"

// TODO: create in menu a NewMenu that receivees titles and items[names, values] and returns a menu
func MenuModel(m []models.ModelEnum, selected models.ModelEnum) string {
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

	return menu.Display(idxSelected)
}

func MenuProfile(profiles map[string]string, selected string) string {
	menu := NewMenu("Select the profile you want to activate:")
	idxSelected := 0
	idx := 0
	for slug, name := range profiles {
		if slug == selected {
			idxSelected = idx
		}
		menu.AddItem(name, slug)
		idx++
	}

	return menu.Display(idxSelected)
}
