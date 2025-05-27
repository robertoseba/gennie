package openai

import (
	"log/slog"
	"net/http"

	"github.com/openai/openai-go/option"
)

func WithModel(model string) opts {
	return func(p *provider) {
		p.model = model
	}
}

func WithHttpClient(client *http.Client) opts {
	return func(p *provider) {
		p.options = append(p.options, option.WithHTTPClient(client))
	}
}

func WithBaseUrl(baseUrl string) opts {
	return func(p *provider) {
		p.options = append(p.options, option.WithBaseURL(baseUrl))
	}
}

func WithLogger(logger *slog.Logger) opts {
	return func(p *provider) {
		p.options = append(p.options, option.WithMiddleware(debugMiddleware(logger)))
		p.logger = logger
	}
}
