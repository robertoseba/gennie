package conversation

import (
	"errors"
)

var ErrAnswerAlreadySet = errors.New("answer already set")

type qa struct {
	Answer   message
	Question message
}
type message struct {
	Content string
	Role    string
}

func NewQA(question string) *qa {
	return &qa{
		Question: message{
			Content: question,
			Role:    userRole,
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
		Content: answer,
		Role:    assistantRole,
	}

	return nil
}
