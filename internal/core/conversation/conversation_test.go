package conversation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversation(t *testing.T) {
	t.Run("NewConversation", func(t *testing.T) {
		c := NewConversation()

		assert.Equal(t, 0, c.Len())
		assert.Equal(t, "", c.LastQuestion())
		assert.Equal(t, "", c.LastAnswer())
		assert.Equal(t, c.CreatedAt, c.UpdatedAt)
	})

	t.Run("Adds a question", func(t *testing.T) {
		c := NewConversation()

		assert.NoError(t, c.NewQuestion("What is your name?"))

		assert.Equal(t, 1, c.Len())
		assert.Equal(t, "What is your name?", c.LastQuestion())
		assert.Equal(t, "", c.LastAnswer())
		assert.True(t, c.CreatedAt.Before(c.UpdatedAt))
	})

	t.Run("Answers question", func(t *testing.T) {
		c := NewConversation()

		assert.NoError(t, c.NewQuestion("What is your name?"))
		assert.NoError(t, c.AnswerLastQuestion("My name is Assistant"))

		assert.Equal(t, 1, c.Len())
		assert.Equal(t, "What is your name?", c.LastQuestion())
		assert.Equal(t, "My name is Assistant", c.LastAnswer())
	})

	t.Run("Adds new question after answering previous one", func(t *testing.T) {
		c := NewConversation()

		assert.NoError(t, c.NewQuestion("What is your name?"))
		assert.NoError(t, c.AnswerLastQuestion("My name is Assistant"))

		assert.NoError(t, c.NewQuestion("What is your age?"))

		assert.Equal(t, 2, c.Len())
		assert.Equal(t, "What is your age?", c.LastQuestion())
		assert.Equal(t, "", c.LastAnswer())

		assert.NoError(t, c.AnswerLastQuestion("I am 1 year old"))
		assert.Equal(t, "I am 1 year old", c.LastAnswer())
	})

	t.Run("Fail when adding new question before answering previous one", func(t *testing.T) {
		c := NewConversation()

		assert.NoError(t, c.NewQuestion("What is your name?"))
		assert.ErrorIs(t, ErrNewQuestionBeforeAnswer, c.NewQuestion("What is your age?"))
	})

	t.Run("Clears conversation", func(t *testing.T) {
		c := NewConversation()
		assert.NoError(t, c.NewQuestion("What is your name?"))
		assert.Equal(t, 1, c.Len())

		c.Clear()

		assert.Equal(t, 0, c.Len())
		assert.Equal(t, "", c.LastQuestion())
	})
}
