package mocks

import (
	"github.com/stretchr/testify/mock"
)

type ClientApiMock struct {
	mock.Mock
}

func (m *ClientApiMock) Post(url string, body string, headers map[string]string) ([]byte, error) {
	args := m.Called(url, body, headers)
	return args.Get(0).([]byte), args.Error(1)
}

func NewClientApiMock() *ClientApiMock {
	return &ClientApiMock{}
}
