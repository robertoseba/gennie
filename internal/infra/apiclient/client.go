package apiclient

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/robertoseba/gennie/internal/core/models"
)

type IhttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	httpClient IhttpClient
}

func NewApiClient(timeout time.Duration) *ApiClient {
	return &ApiClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *ApiClient) Post(url string, body string, headers map[string]string) ([]byte, error) {
	res, err := c.request("POST", url, body, headers)
	if err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) {
			return nil, errors.New("request timeout")
		}
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return nil, errors.New("unauthorized request. Please run 'gennie config' to set your API key.")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unknown error in request. Status code: " + res.Status + "\nBody: " + string(resBody))
	}

	return resBody, nil
}

func (c *ApiClient) PostWithStreaming(url string, body string, headers map[string]string, parser models.ProviderStreamParser) <-chan models.StreamResponse {
	respChan := make(chan models.StreamResponse, 5)

	go func(inputStream chan<- models.StreamResponse) {
		defer close(inputStream)

		res, err := c.request("POST", url, body, headers)
		if err != nil {
			if errors.Is(err, http.ErrHandlerTimeout) {
				inputStream <- models.StreamResponse{Err: errors.New("request timeout")}
				return
			}
			inputStream <- models.StreamResponse{Err: err}
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				inputStream <- models.StreamResponse{Err: err}
				return
			}
			inputStream <- models.StreamResponse{Err: errors.New(string(body))}
			return
		}

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			data := scanner.Bytes()
			parsedData, err := parser(data)
			if err != nil {
				inputStream <- models.StreamResponse{Err: err}
				return
			}
			inputStream <- models.StreamResponse{Data: parsedData, Err: nil}
		}
	}(respChan)

	return respChan
}

func (c *ApiClient) request(method string, url string, body string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))

	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
