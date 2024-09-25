package models

import (
	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/anthropic"
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAIMini ModelEnum = "gpt-4o-mini"
	OpenAI               = "gpt-4o"
	Claude               = "claude-3-5-sonnet"
	Maritaca             = "maritaca"
)

const DefaultModel = OpenAIMini

func (m ModelEnum) String() string {
	switch m {
	case OpenAIMini:
		return "GPT-4o-mini (OPENAI)"
	case OpenAI:
		return "GPT-4o (OPENAI)"
	case Claude:
		return "Claude Sonnet 3.5 (ANTHROPIC)"
	case Maritaca:
		return "Maritaca (USP-BR)"
	default:
		panic("Invalid model")
	}
}

func NewModel(modelType ModelEnum, client *httpclient.HttpClient) IModel {
	model := map[ModelEnum]func(*httpclient.HttpClient) IModel{
		OpenAIMini: func(*httpclient.HttpClient) IModel {
			return openai.NewModel(client, string(modelType))
		},
		OpenAI: func(*httpclient.HttpClient) IModel {
			return openai.NewModel(client, string(modelType))
		},
		Claude: func(*httpclient.HttpClient) IModel {
			return anthropic.NewModel(client, string(modelType))
		},
	}

	activeModel, ok := model[modelType]

	if !ok {
		return model[DefaultModel](client)
	}

	return activeModel(client)
}

func ListModels() []ModelEnum {
	return []ModelEnum{OpenAI, OpenAIMini, Claude, Maritaca}
}
