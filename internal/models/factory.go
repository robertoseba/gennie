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
	model := map[ModelEnum]func(*httpclient.HttpClient) IModel{
		OpenAI: func(*httpclient.HttpClient) IModel {
			return openai.NewModel(client)
		},
		Claude: func(*httpclient.HttpClient) IModel {
			panic("Model not implemented")
		},
	}

	defaultModel := openai.NewModel

	activeModel, ok := model[modelType]

	if !ok {
		return defaultModel(client)
	}

	return activeModel(client)
}

func ListModels() []ModelEnum {
	return []ModelEnum{OpenAI, Claude, Maritaca}
}
