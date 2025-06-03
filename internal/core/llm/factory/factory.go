package factory

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/llm"
	"github.com/robertoseba/gennie/internal/core/llm/anthropic"
	"github.com/robertoseba/gennie/internal/core/llm/gemini"
	"github.com/robertoseba/gennie/internal/core/llm/openai"
)

var ErrModelNotFound = errors.New("model not found")

func NewProvider(modelSlug string, httpClient *http.Client, config config.Config) (llm.LlmProvider, error) {
	modelEnum, _ := ParseFrom(modelSlug)

	switch modelEnum {
	case ClaudeSonnet:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), httpClient), nil
	case Haiku:
		return anthropic.NewProvider(config.APIKeys.AnthropicApiKey, modelEnum.Slug(), httpClient), nil
	case Gemini:
		return gemini.NewProvider(config.APIKeys.GeminiApiKey, modelEnum.Slug(), httpClient), nil
	case OpenAI:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(httpClient),
		), nil
	case OpenAIMini:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(httpClient),
		), nil
	case Groq:
		return openai.NewProvider(config.APIKeys.GroqApiKey,
			openai.WithModel("qwen-qwq-32b"),
			openai.WithHttpClient(httpClient),
			openai.WithBaseUrl("https://api.groq.com/openai/v1"),
		), nil

	case Ollama:
		return openai.NewProvider(config.APIKeys.OpenAiApiKey,
			openai.WithModel(config.Ollama.Model),
			openai.WithHttpClient(httpClient),
			openai.WithBaseUrl(config.Ollama.Host),
		), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrModelNotFound, modelSlug)
}
