package conversation

import (
	"errors"
	"time"
)

var ErrAnswerAlreadySet = errors.New("answer already set")

type qa struct {
	Answer   message
	Question message
	Duration time.Duration
}
type message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

func NewQA(question string) *qa {
	return &qa{
		Question: message{
			Content:   question,
			Role:      userRole,
			Timestamp: time.Now(),
		},
	}
}

func (r *qa) HasAnswer() bool {
	return r.Answer.Content != ""
}

func (r *qa) GetAnswer() string {
	return r.Answer.Content
}

func (r *qa) GetQuestion() string {
	return r.Question.Content
}

func (r *qa) addAnswer(answer string) error {
	if r.HasAnswer() {
		return ErrAnswerAlreadySet
	}

	r.Answer = message{
		Content:   answer,
		Role:      assistantRole,
		Timestamp: time.Now(),
	}
	r.Duration = r.Answer.Timestamp.Sub(r.Question.Timestamp)

	return nil
}

func (r *qa) DurationSeconds() float64 {
	return r.Duration.Seconds()
}
