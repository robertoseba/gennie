package apiclient

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	client := NewApiClient(time.Second * 120)
	assert.NotNil(t, client)
	assert.Equal(t, &http.Client{Timeout: time.Second * 120}, client.httpClient)
}

func TestPost(t *testing.T) {
	t.Run("Returns success response", func(t *testing.T) {
		client := NewApiClient(time.Second * 120)
		mockHttp := NewHttpClientMock()
		client.httpClient = mockHttp

		headers := map[string]string{
			"Content-Type": "application/json",
			"Api-Key":      "Bearer 123",
		}

		mockHttp.On("Do", "http://localhost:8080", "POST", "Bearer 123").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"text": "hi"}`)),
		}, nil)

		body := `{"text": "hello"}`

		resp, err := client.Post("http://localhost:8080", body, headers)

		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"text": "hi"}`), resp)
	})

	t.Run("fails with unauthorized error", func(t *testing.T) {
		client := NewApiClient(time.Second * 120)
		mockHttp := NewHttpClientMock()
		client.httpClient = mockHttp

		mockHttp.On("Do", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"error": "Unauthorized"}`)),
		}, nil)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		body := `{"text": "hello"}`

		resp, err := client.Post("http://localhost:8080", body, headers)

		assert.NotNil(t, err)
		assert.Equal(t, "unauthorized request. Please run 'gennie config' to set your API key.", err.Error())
		assert.Nil(t, resp)
	})

	t.Run("fails with 404 error", func(t *testing.T) {
		client := NewApiClient(time.Second * 120)
		mockHttp := NewHttpClientMock()
		client.httpClient = mockHttp

		mockHttp.On("Do", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Body:       io.NopCloser(strings.NewReader(`{"error": "Not Found"}`)),
		}, nil)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		body := `{"text": "hello"}`

		resp, err := client.Post("http://localhost:8080", body, headers)

		assert.NotNil(t, err)
		assert.Equal(t, "unknown error in request. Status code: 404 Not Found\nBody: {\"error\": \"Not Found\"}", err.Error())
		assert.Nil(t, resp)
	})

	t.Run("fails with timeout", func(t *testing.T) {
		client := NewApiClient(time.Millisecond * 100)
		mockHttp := NewHttpClientMock()
		client.httpClient = mockHttp

		mockHttp.On("Do", mock.Anything, mock.Anything, mock.Anything).Return(&http.Response{}, http.ErrHandlerTimeout)

		headers := map[string]string{
			"Content-Type": "application/json",
		}
		body := `{"text": "hello"}`

		resp, err := client.Post("http://localhost:8080", body, headers)

		assert.NotNil(t, err)
		assert.Equal(t, "request timeout", err.Error())
		assert.Nil(t, resp)
	})
}

type MockedHttpClient struct {
	mock.Mock
}

func (m *MockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req.URL.String(), req.Method, req.Header.Get("Api-Key"))
	return args.Get(0).(*http.Response), args.Error(1)
}

func NewHttpClientMock() *MockedHttpClient {
	return &MockedHttpClient{}

}
