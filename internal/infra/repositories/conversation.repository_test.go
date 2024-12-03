package repositories

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	repo := NewConversationRepository("cacheDir")
	assert.Equal(t, "cacheDir", repo.cacheDir)
}

func TestSaveAndLoadActive(t *testing.T) {
	t.Run("when no active conversation, returns a new one", func(t *testing.T) {
		repo := NewConversationRepository(".")
		conv, err := repo.LoadActive()

		assert.NoError(t, err)
		assert.Equal(t, 0, conv.Len())
	})

	repo := NewConversationRepository("./stub")
	newConversation := conversation.NewConversation("profile-slug", "model-slug")
	err := newConversation.NewQuestion("Question 1")
	assert.NoError(t, err)
	newConversation.AnswerLastQuestion("Answer 1")
	newConversation.NewQuestion("Question 2")
	newConversation.AnswerLastQuestion("Answer 2")

	t.Run("Save active conversation", func(t *testing.T) {
		err = repo.SaveAsActive(newConversation)
		assert.NoError(t, err)

		assert.FileExists(t, "./stub/active.json")
	})

	t.Run("Loaded active conversation should match saved one", func(t *testing.T) {
		loadedConv, err := repo.LoadActive()
		assert.NoError(t, err)

		assert.Equal(t, newConversation.QAs, loadedConv.QAs)
		assert.Equal(t, newConversation.ProfileSlug, loadedConv.ProfileSlug)
		assert.Equal(t, newConversation.ModelSlug, loadedConv.ModelSlug)
		assert.Equal(t, newConversation.CreatedAt.Unix(), loadedConv.CreatedAt.Unix())
	})

	err = os.Remove("./stub/active.json")
	assert.NoError(t, err)
}

func TestExportToFileAndLoadFromFile(t *testing.T) {
	repo := NewConversationRepository("./stub")
	newConversation := conversation.NewConversation("profile-slug", "model-slug")
	err := newConversation.NewQuestion("Question exported 1")
	assert.NoError(t, err)
	newConversation.AnswerLastQuestion("Answer exported 1")

	t.Run("Export conversation to file", func(t *testing.T) {
		err = repo.ExportToFile(newConversation, "./stub/exported_conversation.json")
		assert.NoError(t, err)

		assert.FileExists(t, "./stub/exported_conversation.json")
	})

	t.Run("Loaded conversation should match saved one", func(t *testing.T) {
		loadedConv, err := repo.LoadFromFile("./stub/exported_conversation.json")
		assert.NoError(t, err)

		assert.Equal(t, "Question exported 1", loadedConv.LastQuestion())
		assert.Equal(t, "Answer exported 1", loadedConv.LastAnswer())
		assert.Equal(t, "profile-slug", loadedConv.ProfileSlug)
		assert.Equal(t, "model-slug", loadedConv.ModelSlug)
		assert.Equal(t, newConversation.CreatedAt.Unix(), loadedConv.CreatedAt.Unix())
	})

	err = os.Remove("./stub/exported_conversation.json")
	assert.NoError(t, err)
}
