package factory

import (
	"github.com/robertoseba/gennie/internal/core/llmproviders/anthropic"
	"github.com/robertoseba/gennie/internal/core/llmproviders/gemini"
)

type ModelEnum string

const (
	ClaudeSonnet ModelEnum = anthropic.ExportedSonnetSlug
	Haiku        ModelEnum = anthropic.ExportedHaikuSlug
	Gemini       ModelEnum = gemini.ExportedModelSlug
	OpenAIMini   ModelEnum = "gpt-4.1-mini"
	OpenAI       ModelEnum = "gpt-4.1"
	Groq         ModelEnum = "groq"
	Ollama       ModelEnum = "ollama"
)

const DefaultModel = ClaudeSonnet

var availableModels = map[ModelEnum]string{
	ClaudeSonnet: anthropic.ExportedSonnetDescription,
	Haiku:        anthropic.ExportedHaikuDescription,
	Gemini:       gemini.ExportedModelDescription,
	OpenAIMini:   "GPT-4.1-mini (OPENAI)",
	OpenAI:       "GPT-4.1 (OPENAI)",
	Groq:         "Groq (DeepSeek-R1-Distill-Llama-70B)",
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
		return DefaultModel, false
	}

	return ModelEnum(modelSlug), true
}

func ListModels() map[ModelEnum]string {
	return availableModels
}
