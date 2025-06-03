package factory

import (
	"github.com/robertoseba/gennie/internal/core/llm/anthropic"
	"github.com/robertoseba/gennie/internal/core/llm/gemini"
	"github.com/robertoseba/gennie/internal/core/llm/openai"
)

type ModelEnum string

const (
	ClaudeSonnet ModelEnum = anthropic.ExportedSonnetSlug
	Haiku        ModelEnum = anthropic.ExportedHaikuSlug
	Gemini       ModelEnum = gemini.ExportedModelSlug
	OpenAIMini   ModelEnum = openai.ExportedModelSlug_Mini
	OpenAI       ModelEnum = openai.ExportedModelSlug
	Groq         ModelEnum = "groq"
	Ollama       ModelEnum = "ollama"
)

const DefaultModel = ClaudeSonnet

var availableModels = map[ModelEnum]string{
	ClaudeSonnet: anthropic.ExportedSonnetDescription,
	Haiku:        anthropic.ExportedHaikuDescription,
	Gemini:       gemini.ExportedModelDescription,
	OpenAIMini:   openai.ExportedModelDescription_Mini,
	OpenAI:       openai.ExportedModelDescription,
	Groq:         "Groq (Qwen-QWQ-32B)",
	Ollama:       "Ollama",
}

func (m ModelEnum) String() string {
	return availableModels[m]
}

func (m ModelEnum) Slug() string {
	return string(m)
}

func ParseFrom(modelSlug string) (ModelEnum, bool) {
	_, ok := availableModels[ModelEnum(modelSlug)]
	if !ok {
		return "", false
	}

	return ModelEnum(modelSlug), true
}

func ListModels() map[ModelEnum]string {
	return availableModels
}
