package chat

import "time"

const UserRole = "user"
const AssistantRole = "assistant"
const SystemRole = "system"

// Data sent to the model
type input struct {
	content   string
	role      string
	timestamp time.Time
}

// Data returned by the model
type output struct {
	content   string
	role      string
	timestamp time.Time
}

type Response struct {
	answer   output
	question input
	duration time.Duration
}

func (r *Response) Answer() string {
	return r.answer.content
}

func (r *Response) AddQuestion(question string) {
	r.question = input{
		content:   question,
		role:      UserRole,
		timestamp: time.Now(),
	}
}

func (r *Response) AddAnswer(answer string) {
	r.answer = output{
		content:   answer,
		role:      AssistantRole,
		timestamp: time.Now(),
	}
	r.duration = r.answer.timestamp.Sub(r.question.timestamp)
}

func (r *Response) DurationSeconds() float64 {
	return r.duration.Seconds()
}
