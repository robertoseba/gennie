package conversation

import (
	"errors"
	"testing"
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

		if qa.HasAnswer() {
			t.Errorf("Expected HasAnswer to be false, got true")
		}

		if qa.GetAnswer() != "" {
			t.Errorf("Expected answer to be '', got %s", qa.GetAnswer())
		}
	})

	t.Run("has already filled", func(t *testing.T) {
		qa := NewQA("question")
		qa.addAnswer("answer")
		if qa.GetAnswer() != "answer" {
			t.Errorf("Expected answer to be 'answer', got %s", qa.GetAnswer())
		}

		if !qa.HasAnswer() {
			t.Errorf("Expected HasAnswer to be true, got false")
		}
	})

	t.Run("Can change answer already set", func(t *testing.T) {
		qa := NewQA("question")
		qa.addAnswer("answer")
		err := qa.addAnswer("answer2")

		if !errors.Is(err, ErrAnswerAlreadySet) {
			t.Errorf("Expected error to be ErrAnswerAlreadySet, got %v", err)
		}
	})
}
