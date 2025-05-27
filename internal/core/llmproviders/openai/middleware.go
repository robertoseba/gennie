package openai

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/openai/openai-go/option"
)

func debugMiddleware(logger *slog.Logger) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		reqLogger := logger.WithGroup("Request").With("URL", r.URL.String()).With("Method", r.Method)
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				reqLogger.Error("Failed to read request body", "error", err)
			}
			reqLogger.Debug("Sending request", "body", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the body for the next handler
		}
		n, err := next(r)
		if n != nil {
			respLogger := logger.WithGroup("Response").With("StatusCode", n.StatusCode)
			if n.Body != nil {
				body, err := io.ReadAll(n.Body)
				if err != nil {
					respLogger.Error("Failed to read response body", "error", err)
				}
				respLogger.Debug("Received Response", "body", string(body))
				n.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the body for the next handler
			}
		}
		return n, err
	}
}
