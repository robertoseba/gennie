package base

import "errors"

var (
	ErrEmptyConversation           = errors.New("there are no questions to answer")
	ErrLastQuestionAlreadyAnswered = errors.New("last conversation has already been answered")
	ErrModelNotFound               = errors.New("model not found")
)
