package usecases

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/infra/repositories"
	"github.com/stretchr/testify/require"
)

func TestExportConversationService_Execute(t *testing.T) {
	t.Run("Loads active conversation and saves to a file", func(t *testing.T) {
		//Tests for file content are being done in the repository test
		conv, repo := setupActiveConversation(t)
		repo.ExportToFile(conv, "./conversation.json")

		require.FileExists(t, "./conversation.json")
		require.NoError(t, os.Remove("./conversation.json"))
	})
}

func setupActiveConversation(t *testing.T) (*conversation.Conversation, conversation.IConversationRepository) {
	repo := repositories.NewConversationRepository(".")
	conv, err := repo.LoadActive()
	require.NoError(t, err)
	conv.NewQuestion("What is your name?")
	conv.AnswerLastQuestion("My name is Gennie")
	return conv, repo
}
