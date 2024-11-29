package conversation

import (
	"errors"
	"testing"
)

func TestConversation(t *testing.T) {
	t.Run("NewConversation", func(t *testing.T) {
		c := NewConversation()
		if c.Len() != 0 {
			t.Errorf("Expected length to be 0, got %d", c.Len())
		}

		if c.LastQuestion() != "" {
			t.Errorf("Expected last question to be '', got %s", c.LastQuestion())
		}

		if c.LastAnswer() != "" {
			t.Errorf("Expected last answer to be '', got %s", c.LastAnswer())
		}
	})

	t.Run("Adds a question", func(t *testing.T) {
		c := NewConversation()
		if err := c.NewQuestion("What is your name?"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if c.Len() != 1 {
			t.Errorf("Expected length to be 1, got %d", c.Len())
		}

		if c.LastQuestion() != "What is your name?" {
			t.Errorf("Expected last question to be 'What is your name?', got %s", c.LastQuestion())
		}

		if c.LastAnswer() != "" {
			t.Errorf("Expected last answer to be '', got %s", c.LastAnswer())
		}
	})

	t.Run("Answers question", func(t *testing.T) {
		c := NewConversation()
		if err := c.NewQuestion("What is your name?"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if err := c.AnswerLastQuestion("My name is Assistant"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if c.LastAnswer() != "My name is Assistant" {
			t.Errorf("Expected last answer to be 'My name is Assistant', got %s", c.LastAnswer())
		}
	})

	t.Run("Adds new question after answering previous one", func(t *testing.T) {
		c := NewConversation()
		if err := c.NewQuestion("What is your name?"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if err := c.AnswerLastQuestion("My name is Assistant"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if err := c.NewQuestion("What is your age?"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if c.Len() != 2 {
			t.Errorf("Expected length to be 2, got %d", c.Len())
		}

		if c.LastQuestion() != "What is your age?" {
			t.Errorf("Expected last question to be 'What is your age?', got %s", c.LastQuestion())
		}

		if c.LastAnswer() != "" {
			t.Errorf("Expected last answer to be '', got %s", c.LastAnswer())
		}

		if err := c.AnswerLastQuestion("I am 1 year old"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if c.LastAnswer() != "I am 1 year old" {
			t.Errorf("Expected last answer to be 'I am 1 year old', got %s", c.LastAnswer())
		}
	})

	t.Run("Fail when adding new question before answering previous one", func(t *testing.T) {
		c := NewConversation()
		if err := c.NewQuestion("What is your name?"); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		err := c.NewQuestion("What is your age?")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if errors.Is(err, ErrNewQuestionBeforeAnswer) == false {
			t.Errorf("Expected error to be ErrNewQuestionBeforeAnswer, got %v", err)
		}
	})
}
