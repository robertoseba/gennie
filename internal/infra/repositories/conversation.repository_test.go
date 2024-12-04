package repositories

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	repo := NewConversationRepository("cacheDir")
	require.Equal(t, "cacheDir", repo.cacheDir)
}

func TestSaveAndLoadActive(t *testing.T) {
	t.Run("when no active conversation, returns a new one", func(t *testing.T) {
		repo := NewConversationRepository(".")
		conv, err := repo.LoadActive()

		require.NoError(t, err)
		require.Equal(t, 0, conv.Len())
	})

	repo := NewConversationRepository("./stub")
	newConversation := conversation.NewConversation("profile-slug", "model-slug")
	err := newConversation.NewQuestion("Question 1")
	require.NoError(t, err)
	newConversation.AnswerLastQuestion("Answer 1")
	newConversation.NewQuestion("Question 2")
	newConversation.AnswerLastQuestion("Answer 2")

	t.Run("Save active conversation", func(t *testing.T) {
		err = repo.SaveAsActive(newConversation)
		require.NoError(t, err)

		require.FileExists(t, "./stub/active.json")
	})

	t.Run("Loaded active conversation should match saved one", func(t *testing.T) {
		loadedConv, err := repo.LoadActive()
		require.NoError(t, err)

		require.Equal(t, newConversation.QAs, loadedConv.QAs)
		require.Equal(t, newConversation.ProfileSlug, loadedConv.ProfileSlug)
		require.Equal(t, newConversation.ModelSlug, loadedConv.ModelSlug)
		require.Equal(t, newConversation.CreatedAt.Unix(), loadedConv.CreatedAt.Unix())
	})

	err = os.Remove("./stub/active.json")
	require.NoError(t, err)
}

func TestExportToFileAndLoadFromFile(t *testing.T) {
	repo := NewConversationRepository("./stub")
	newConversation := conversation.NewConversation("profile-slug", "model-slug")
	err := newConversation.NewQuestion("Question exported 1")
	require.NoError(t, err)
	newConversation.AnswerLastQuestion("Answer exported 1")

	t.Run("Export conversation to file", func(t *testing.T) {
		err = repo.ExportToFile(newConversation, "./stub/exported_conversation.json")
		require.NoError(t, err)

		require.FileExists(t, "./stub/exported_conversation.json")
	})

	t.Run("Loaded conversation should match saved one", func(t *testing.T) {
		loadedConv, err := repo.LoadFromFile("./stub/exported_conversation.json")
		require.NoError(t, err)

		require.Equal(t, "Question exported 1", loadedConv.LastQuestion())
		require.Equal(t, "Answer exported 1", loadedConv.LastAnswer())
		require.Equal(t, "profile-slug", loadedConv.ProfileSlug)
		require.Equal(t, "model-slug", loadedConv.ModelSlug)
		require.Equal(t, newConversation.CreatedAt.Unix(), loadedConv.CreatedAt.Unix())
	})

	err = os.Remove("./stub/exported_conversation.json")
	require.NoError(t, err)
}
