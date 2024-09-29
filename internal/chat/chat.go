package chat

import "time"

const UserRole = "user"
const AssistantRole = "assistant"
const SystemRole = "system"

// Data sent to the model
type input struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// Data returned by the model
type output struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type Chat struct {
	Answer   output
	Question input
	Duration time.Duration
}

func NewChat(question string) *Chat {
	return &Chat{
		Question: input{
			Content:   question,
			Role:      UserRole,
			Timestamp: time.Now(),
		},
	}
}

func (r *Chat) HasAnswer() bool {
	return r.Answer.Content != ""
}

func (r *Chat) GetAnswer() string {
	return r.Answer.Content
}

func (r *Chat) GetQuestion() string {
	return r.Question.Content
}

func (r *Chat) AddQuestion(question string) {
	r.Question = input{
		Content:   question,
		Role:      UserRole,
		Timestamp: time.Now(),
	}
}

func (r *Chat) AddAnswer(answer string) {
	r.Answer = output{
		Content:   answer,
		Role:      AssistantRole,
		Timestamp: time.Now(),
	}
	r.Duration = r.Answer.Timestamp.Sub(r.Question.Timestamp)
}

func (r *Chat) DurationSeconds() float64 {
	return r.Duration.Seconds()
}
