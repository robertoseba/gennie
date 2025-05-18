package output

import "github.com/robertoseba/gennie/internal/core/llmproviders/factory"

// TODO: create in menu a NewMenu that receives titles and items[names, values] and returns a menu
func MenuModel(m map[factory.ModelEnum]string, selected factory.ModelEnum) factory.ModelEnum {
	menu := NewMenu("Select a model:")

	idxSelected := 0
	idx := 0
	for model, desc := range m {
		if model == selected {
			idxSelected = idx
		}
		menu.AddItem(desc, model.Slug())
		idx++
	}

	return factory.ModelEnum(menu.Display(idxSelected))
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
