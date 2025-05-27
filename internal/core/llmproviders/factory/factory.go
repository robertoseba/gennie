package factory

import (
	"log/slog"
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"github.com/robertoseba/gennie/internal/core/llmproviders/anthropic"
	"github.com/robertoseba/gennie/internal/core/llmproviders/gemini"
	"github.com/robertoseba/gennie/internal/core/llmproviders/openai"
)

func NewProvider(modelEnum ModelEnum, logger *slog.Logger, httpClient *http.Client, config config.Config) llmcore.LlmProvider {
	switch modelEnum {
	case ClaudeSonnet:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), logger, httpClient)
	case Haiku:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), logger, httpClient)
	case Gemini:
		return gemini.NewProvider(config.APIKeys.GeminiApiKey, modelEnum.Slug(), httpClient)
	case OpenAI:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(httpClient),
			openai.WithLogger(logger),
		)
	case OpenAIMini:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(httpClient),
			openai.WithLogger(logger),
		)
	case Groq:
		return openai.NewProvider(config.APIKeys.GroqApiKey,
			openai.WithModel("qwen-qwq-32b"),
			openai.WithHttpClient(httpClient),
			openai.WithLogger(logger),
			openai.WithBaseUrl("https://api.groq.com/openai/v1"),
		)

	case Ollama:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(config.Ollama.Model),
			openai.WithHttpClient(httpClient),
			openai.WithLogger(logger),
			openai.WithBaseUrl(config.Ollama.Host),
		)
	default:
		panic("Unsupported model: " + modelEnum.Slug())
	}
}
