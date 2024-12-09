package mocks

import (
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/stretchr/testify/mock"
)

type ClientApiMock struct {
	mock.Mock
}

func (m *ClientApiMock) Post(url string, body string, headers map[string]string) ([]byte, error) {
	args := m.Called(url, body, headers)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *ClientApiMock) PostWithStreaming(url string,
	body string,
	headers map[string]string,
	parser models.ProviderStreamParser) <-chan models.StreamResponse {

	args := m.Called(url, body, headers, parser)
	return args.Get(0).(<-chan models.StreamResponse)
}

func NewClientApiMock() *ClientApiMock {
	return &ClientApiMock{}
}
