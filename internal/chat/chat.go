package chat

import (
	"errors"
	"time"
)

const UserRole = "user"
const AssistantRole = "assistant"
const SystemRole = "system"

var ErrAnswerAlreadySet = errors.New("answer already set")

// Data sent to the model
type message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type QA struct {
	Answer   message
	Question message
	Duration time.Duration
}

func NewQA(question string) *QA {
	return &QA{
		Question: message{
			Content:   question,
			Role:      UserRole,
			Timestamp: time.Now(),
		},
	}
}

func (r *QA) HasAnswer() bool {
	return r.Answer.Content != ""
}

func (r *QA) GetAnswer() string {
	return r.Answer.Content
}

func (r *QA) GetQuestion() string {
	return r.Question.Content
}

func (r *QA) AddAnswer(answer string) error {
	if r.HasAnswer() {
		return ErrAnswerAlreadySet
	}

	r.Answer = message{
		Content:   answer,
		Role:      AssistantRole,
		Timestamp: time.Now(),
	}
	r.Duration = r.Answer.Timestamp.Sub(r.Question.Timestamp)

	return nil
}

func (r *QA) DurationSeconds() float64 {
	return r.Duration.Seconds()
}
