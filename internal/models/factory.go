package models

import (
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAI   ModelEnum = "OpenAI"
	Claude             = "Claude"
	Maritaca           = "Maritaca"
)

func NewModel(modelType ModelEnum) IModel {
	switch modelType {
	case OpenAI:
		return openai.NewModel()

	// case Claude:
	// 	return NewClaude()
	// case Maritaca:
	// 	return NewMaritaca()
	default:
		return openai.NewModel()
	}
}
