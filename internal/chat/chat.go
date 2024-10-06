package chat

import "time"

const UserRole = "user"
const AssistantRole = "assistant"
const SystemRole = "system"

// Data sent to the model
type message struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type Chat struct {
	Answer   message
	Question message
	Duration time.Duration
}

func NewChat(question string) *Chat {
	return &Chat{
		Question: message{
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
	r.Question = message{
		Content:   question,
		Role:      UserRole,
		Timestamp: time.Now(),
	}
}

func (r *Chat) AddAnswer(answer string) {
	r.Answer = message{
		Content:   answer,
		Role:      AssistantRole,
		Timestamp: time.Now(),
	}
	r.Duration = r.Answer.Timestamp.Sub(r.Question.Timestamp)
}

func (r *Chat) DurationSeconds() float64 {
	return r.Duration.Seconds()
}
