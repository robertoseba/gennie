package mocks

import (
	"github.com/robertoseba/gennie/internal/core/profile"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func NewMockProfileRepository() *MockProfileRepository {
	return &MockProfileRepository{}
}

func (m *MockProfileRepository) FindBySlug(slug string) (*profile.Profile, error) {
	args := m.Called(slug)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*profile.Profile), args.Error(1)
}

func (m *MockProfileRepository) ListAll() (map[string]*profile.Profile, error) {
	args := m.Called()

	return args.Get(0).(map[string]*profile.Profile), args.Error(1)
}
