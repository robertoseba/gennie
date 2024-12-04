package usecases

import (
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/profile"
	"github.com/robertoseba/gennie/internal/infra/repositories"
	"github.com/robertoseba/gennie/internal/infra/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileListAll(t *testing.T) {
	profRepo := repositories.NewProfileRepository(".")
	mockConvRepo := mocks.NewMockConversationRepository()

	service := NewSelectProfileService(profRepo, mockConvRepo)
	p, err := service.ListAll()

	require.NoError(t, err)
	require.Len(t, p, 1)
	require.Equal(t, profile.DefaultProfile(), p["default"])
}

func TestProfileSetAsActive(t *testing.T) {
	t.Run("sets profile as active", func(t *testing.T) {
		activeConv := conversation.NewConversation("profile-slug", "model22")

		profRepo := repositories.NewProfileRepository(".")
		mockConvRepo := mocks.NewMockConversationRepository()
		mockConvRepo.On("LoadActive").Return(activeConv, nil)
		mockConvRepo.On("SaveAsActive", activeConv).Return(nil)

		assert.Equal(t, "profile-slug", activeConv.ProfileSlug)

		activeProfile := profile.DefaultProfile()

		service := NewSelectProfileService(profRepo, mockConvRepo)
		err := service.SetAsActive(activeProfile)

		require.NoError(t, err)
		require.Equal(t, activeProfile.Slug, activeConv.ProfileSlug)
	})

	t.Run("returns error if cant save as active", func(t *testing.T) {
		activeConv := conversation.NewConversation("profile-slug", "model22")

		profRepo := repositories.NewProfileRepository(".")
		mockConvRepo := mocks.NewMockConversationRepository()
		mockConvRepo.On("LoadActive").Return(activeConv, nil)
		mockConvRepo.On("SaveAsActive", activeConv).Return(assert.AnError)

		activeProfile := profile.DefaultProfile()

		service := NewSelectProfileService(profRepo, mockConvRepo)
		err := service.SetAsActive(activeProfile)

		require.Error(t, err)
	})
}
