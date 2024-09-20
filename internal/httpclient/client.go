package httpclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	timeout     time.Duration
	bearerToken string
}

func NewClient() *HttpClient {
	return &HttpClient{
		timeout:     10,
		bearerToken: "",
	}
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	res, err := c.request("GET", url, "")

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

	return resBody, nil
}

func (c *HttpClient) SetTimeout(timeout int) {
	c.timeout = time.Duration(timeout)
}

func (c *HttpClient) SetBearerToken(authKey string) {
	c.bearerToken = authKey
}

func (c *HttpClient) request(method string, url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
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
