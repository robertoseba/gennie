package mocks

import (
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/mock"
)

type MockConversationRepository struct {
	mock.Mock
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{}
}

func (m *MockConversationRepository) LoadActive() (*conversation.Conversation, error) {
	args := m.Called()
	return args.Get(0).(*conversation.Conversation), args.Error(1)
}

func (m *MockConversationRepository) SaveAsActive(conv *conversation.Conversation) error {
	args := m.Called(conv)
	return args.Error(0)
}
func (m *MockConversationRepository) ExportToFile(conv *conversation.Conversation, filename string) error {
	args := m.Called(conv, filename)
	return args.Error(0)
}
func (m *MockConversationRepository) LoadFromFile(filename string) (*conversation.Conversation, error) {
	args := m.Called(filename)
	return args.Get(0).(*conversation.Conversation), args.Error(1)
}
