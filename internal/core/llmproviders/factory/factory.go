package factory

import (
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"github.com/robertoseba/gennie/internal/core/llmproviders/anthropic"
	"github.com/robertoseba/gennie/internal/core/llmproviders/gemini"
	"github.com/robertoseba/gennie/internal/core/llmproviders/openai"
)

func NewProvider(modelEnum ModelEnum, httpClient *http.Client, config config.Config) llmcore.LlmProvider {
	switch modelEnum {
	case ClaudeSonnet:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), httpClient)
	case Gemini:
		return gemini.NewProvider(config.APIKeys.GeminiApiKey, modelEnum.Slug(), httpClient)
	case OpenAI:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey, modelEnum.Slug(), httpClient)
	case OpenAIMini:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey, modelEnum.Slug(), httpClient)
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
