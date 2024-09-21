package models

import (
	"errors"
	"strings"

	"github.com/robertoseba/gennie/internal/httpclient"
	"github.com/robertoseba/gennie/internal/models/openai"
)

type ModelEnum string

const (
	OpenAI   ModelEnum = "openai"
	Claude             = "claude"
	Maritaca           = "maritaca"
)

func (m *ModelEnum) String() string {
	return string(*m)
}
func (m *ModelEnum) From(s string) (ModelEnum, error) {
	switch strings.ToLower(s) {
	case "openai":
		return OpenAI, nil
	case "claude":
		return Claude, nil
	case "maritaca":
		return Maritaca, nil
	default:
		return "", errors.New("Invalid model")
	}
}

func (m *ModelEnum) All() []ModelEnum {
	return []ModelEnum{OpenAI, Claude, Maritaca}
}

var models = map[ModelEnum]func(*httpclient.HttpClient) IModel{
	OpenAI: func(client *httpclient.HttpClient) IModel {
		return openai.NewModel(client)
	},
	Claude: func(client *httpclient.HttpClient) IModel {
		return nil
	},
}

func NewModel(modelType ModelEnum, client *httpclient.HttpClient) IModel {
	return models[modelType](client)
}
