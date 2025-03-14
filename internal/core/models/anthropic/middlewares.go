package anthropic

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/anthropics/anthropic-sdk-go/option"
)

func NewErrorMiddleware(output io.Writer) option.Middleware {
	return func(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
		// TODO: think of a better way to have a debug flag
		if os.Getenv("DEBUG") == "true" {
			reqLog, err := httputil.DumpRequestOut(req, true)
			if err != nil {
				fmt.Fprintf(output, "Error dumping request: %s\n", err)
			}
			fmt.Fprintf(output, "Request: %s\n\n", string(reqLog))
		}

		res, err = next(req)

		if res.StatusCode != http.StatusOK || err != nil {
			fmt.Fprintf(output, "Response: %s %s\n", res.Status, req.URL)
			var respBody []byte
			if res != nil {
				respBody, err = io.ReadAll(res.Body)
			}
			fmt.Fprintf(output, "Response Body: %s\n", string(respBody))
			if err != nil {
				fmt.Fprintf(output, "Error: %s\n", err)
			}
		}

		return res, err
	}
}
