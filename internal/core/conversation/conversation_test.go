package conversation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConversation(t *testing.T) {
	t.Run("NewConversation", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.Equal(t, 0, c.Len())
		require.Equal(t, "", c.LastQuestion())
		require.Equal(t, "", c.LastAnswer())
		require.Equal(t, c.CreatedAt, c.UpdatedAt)
		require.Equal(t, "profile-slug", c.ProfileSlug)
		require.Equal(t, "model-slug", c.ModelSlug)
	})

	t.Run("Adds a question", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.NoError(t, c.NewQuestion("What is your name?"))

		require.Equal(t, 1, c.Len())
		require.Equal(t, "What is your name?", c.LastQuestion())
		require.Equal(t, "", c.LastAnswer())
		require.True(t, c.CreatedAt.Before(c.UpdatedAt))
	})

	t.Run("Answers question", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.NoError(t, c.NewQuestion("What is your name?"))
		require.NoError(t, c.AnswerLastQuestion("My name is Assistant"))

		require.Equal(t, 1, c.Len())
		require.Equal(t, "What is your name?", c.LastQuestion())
		require.Equal(t, "My name is Assistant", c.LastAnswer())
	})

	t.Run("Adds new question after answering previous one", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.NoError(t, c.NewQuestion("What is your name?"))
		require.NoError(t, c.AnswerLastQuestion("My name is Assistant"))

		require.NoError(t, c.NewQuestion("What is your age?"))

		require.Equal(t, 2, c.Len())
		require.Equal(t, "What is your age?", c.LastQuestion())
		require.Equal(t, "", c.LastAnswer())

		require.NoError(t, c.AnswerLastQuestion("I am 1 year old"))
		require.Equal(t, "I am 1 year old", c.LastAnswer())
	})

	t.Run("Fail when adding new question before answering previous one", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.NoError(t, c.NewQuestion("What is your name?"))
		require.ErrorIs(t, ErrNewQuestionBeforeAnswer, c.NewQuestion("What is your age?"))
	})

	t.Run("Clears conversation", func(t *testing.T) {
		c := NewConversation("profile-slug", "model-slug")

		require.NoError(t, c.NewQuestion("What is your name?"))
		require.Equal(t, 1, c.Len())

		c.Clear()

		require.Equal(t, 0, c.Len())
		require.Equal(t, "", c.LastQuestion())
	})
}
