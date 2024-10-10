package httpclient

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type IHttpClient interface {
	Post(url string, body string, headers map[string]string) ([]byte, error)
}
type HttpClient struct {
	timeout time.Duration
}

func NewClient() *HttpClient {
	return &HttpClient{
		timeout: 40,
	}
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	res, err := c.request("GET", url, "", nil)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) Post(url string, body string, headers map[string]string) ([]byte, error) {
	res, err := c.request("POST", url, body, headers)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return nil, errors.New("Unauthorized request. Please run 'gennie config' to set your API key.")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Unknown error in request. Status code: " + res.Status + "\nBody: " + string(resBody))
	}

	return resBody, nil
}

func (c *HttpClient) SetTimeout(timeout int) {
	c.timeout = time.Duration(timeout)
}

func (c *HttpClient) request(method string, url string, body string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))

	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: c.timeout * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
