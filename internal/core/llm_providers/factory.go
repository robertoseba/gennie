package llmproviders

import (
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llm_providers/anthropic"
	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
)

func NewModel(modelEnum base.ModelEnum, httpClient *http.Client, config config.Config) LlmProvider {
	switch modelEnum {
	// case OpenAI:
	// 	return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	// case OpenAIMini:
	// 	return openai.NewProvider(m.Slug(), config.APIKeys.OpenAiApiKey)
	case base.ClaudeSonnet:
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
