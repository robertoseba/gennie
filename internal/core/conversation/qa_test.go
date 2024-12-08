package conversation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQA(t *testing.T) {
	t.Run("has question", func(t *testing.T) {
		qa := NewQA("question")
		if qa.GetQuestion() != "question" {
			t.Errorf("Expected question to be 'question', got %s", qa.GetQuestion())
		}
	})

	t.Run("answer is empty", func(t *testing.T) {
		qa := NewQA("question")

		require.False(t, qa.HasAnswer())
		require.Equal(t, "", qa.GetAnswer())
	})

	t.Run("has already filled", func(t *testing.T) {
		qa := NewQA("question")
		qa.addAnswer("answer")

		require.True(t, qa.HasAnswer())
		require.Equal(t, "answer", qa.GetAnswer())
	})

	t.Run("Can change answer already set", func(t *testing.T) {
		qa := NewQA("question")
		qa.addAnswer("answer")
		err := qa.addAnswer("answer2")

		require.ErrorIs(t, err, ErrAnswerAlreadySet)
	})

	t.Run("Roles are assigned correctly", func(t *testing.T) {
		qa := NewQA("question")
		qa.addAnswer("answer")

		require.Equal(t, userRole, qa.Question.Role)
		require.Equal(t, assistantRole, qa.Answer.Role)
	})
}
