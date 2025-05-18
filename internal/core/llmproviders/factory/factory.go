package factory

import (
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"github.com/robertoseba/gennie/internal/core/llmproviders/anthropic"
)

func NewModel(modelEnum ModelEnum, httpClient *http.Client, config config.Config) llmcore.LlmProvider {
	switch modelEnum {
	// case OpenAI:
	// 	return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	// case OpenAIMini:
	// 	return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	case ClaudeSonnet:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), httpClient)
	// case Maritaca:
	// 	return maritaca.NewProvider(m.Slug(), config.APIKeys.MaritacaApiKey)
	// case Groq:
	// 	return groq.NewProvider(m.Slug(), config.APIKeys.GroqApiKey)
	// case Ollama:
	// 	return ollama.NewProvider(m.Slug(), config.Ollama.Host, config.Ollama.Model)
	default:
		panic("Unsupported model: " + modelEnum.Slug())
	}
}
