package httpclient

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	timeout    time.Duration
	authHeader string
}

func NewClient() *HttpClient {
	return &HttpClient{
		timeout:    10,
		authHeader: "",
	}
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	res, err := c.request("GET", url, "")

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

func (c *HttpClient) Post(url string, body string) ([]byte, error) {
	res, err := c.request("POST", url, body)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return nil, errors.New("Unauthorized request. Please make sure you have set the correct API key in your environment variables")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Unknown error in request. Status code: " + res.Status + "\nBody: " + string(resBody))
	}

	return resBody, nil
}

func (c *HttpClient) SetTimeout(timeout int) {
	c.timeout = time.Duration(timeout)
}

func (c *HttpClient) SetAuthHeader(authHeader string) {
	c.authHeader = authHeader
}

func (c *HttpClient) request(method string, url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	client := &http.Client{
		Timeout: c.timeout * time.Second,
	}

	//TODO: Implement retry logic
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// TODO: retry in case of specific status code
	// if res.StatusCode != http.StatusOK {
	// 	return nil, &HttpError{
	// 		StatusCode: res.StatusCode,
	// 		Status:     res.Status
	// 	}
	// }

	return res, nil
}
