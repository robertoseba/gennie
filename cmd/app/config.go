package app

import (
	"github.com/robertoseba/gennie/internal/models"
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
