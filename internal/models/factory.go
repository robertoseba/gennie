package models

import (
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAI   ModelEnum = "OpenAI"
	Claude             = "Claude"
	Maritaca           = "Maritaca"
)

func NewModel(modelType ModelEnum, client *httpclient.HttpClient) IModel {
	return openai.NewModel(client)
	// switch modelType {
	// case OpenAI:
	// 	return openai.NewModel(client)

	// // case Claude:
	// // 	return NewClaude()
	// // case Maritaca:
	// // 	return NewMaritaca()
	// default:
	// 	return openai.NewModel(client)
	// }
}
