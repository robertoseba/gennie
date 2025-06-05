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

type ProviderFactory interface {
	CreateProvider(modelSlug string) (llm.Provider, error)
}

type providerFactory struct {
	httpClient *http.Client
	config     config.Config
}

func NewProviderFactory(httpClient *http.Client, config config.Config) *providerFactory {
	return &providerFactory{
		httpClient: httpClient,
		config:     config,
	}
}

func (f *providerFactory) CreateProvider(modelSlug string) (llm.Provider, error) {
	modelEnum, _ := ParseFrom(modelSlug)

	switch modelEnum {
	case ClaudeSonnet:
		return anthropic.NewProvider(f.config.APIKeys.AnthropicApiKey, modelEnum.Slug(), f.httpClient), nil
	case Haiku:
		return anthropic.NewProvider(f.config.APIKeys.AnthropicApiKey, modelEnum.Slug(), f.httpClient), nil
	case Gemini:
		return gemini.NewProvider(f.config.APIKeys.GeminiApiKey, modelEnum.Slug(), f.httpClient), nil
	case OpenAI:
		return openai.NewProvider(f.config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(f.httpClient),
		), nil
	case OpenAIMini:
		return openai.NewProvider(f.config.APIKeys.OpenAiApiKey,
			openai.WithModel(modelEnum.Slug()),
			openai.WithHttpClient(f.httpClient),
		), nil
	case Groq:
		return openai.NewProvider(f.config.APIKeys.GroqApiKey,
			openai.WithModel("qwen-qwq-32b"),
			openai.WithHttpClient(f.httpClient),
			openai.WithBaseUrl("https://api.groq.com/openai/v1"),
		), nil

	case Ollama:
		return openai.NewProvider(f.config.APIKeys.OpenAiApiKey,
			openai.WithModel(f.config.Ollama.Model),
			openai.WithHttpClient(f.httpClient),
			openai.WithBaseUrl(f.config.Ollama.Host),
		), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrModelNotFound, modelSlug)
}
