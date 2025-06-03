package openai

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/openai/openai-go/option"
)

func DebugMiddleware(logger *slog.Logger) option.Middleware {
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

		res, err := next(r)
		if err != nil {
			logger.Error("Response failed", "error", err)
		}

		logger.Debug("Response Status: " + res.Status)
		res.Body = io.NopCloser(io.TeeReader(res.Body, &logStream{logger: logger}))
		return res, err
	}
}

type logStream struct {
	logger *slog.Logger
}

func (l *logStream) Write(p []byte) (n int, err error) {
	l.logger.Debug("received streamed data:", "payload", string(p))
	return len(p), nil
}
