package conversation

import (
	"errors"
	"time"
)

var ErrAnswerAlreadySet = errors.New("answer already set")

// Data sent to the model
type message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type qa struct {
	answer   message
	question message
	duration time.Duration
}

func NewQA(question string) *qa {
	return &qa{
		question: message{
			Content:   question,
			Role:      userRole,
			Timestamp: time.Now(),
		},
	}
}

func (r *qa) HasAnswer() bool {
	return r.answer.Content != ""
}

func (r *qa) GetAnswer() string {
	return r.answer.Content
}

func (r *qa) GetQuestion() string {
	return r.question.Content
}

func (r *qa) AddAnswer(answer string) error {
	if r.HasAnswer() {
		return ErrAnswerAlreadySet
	}

	r.answer = message{
		Content:   answer,
		Role:      assistantRole,
		Timestamp: time.Now(),
	}
	r.duration = r.answer.Timestamp.Sub(r.question.Timestamp)

	return nil
}

func (r *qa) DurationSeconds() float64 {
	return r.duration.Seconds()
}
