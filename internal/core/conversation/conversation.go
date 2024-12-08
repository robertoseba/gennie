package conversation

import (
	"errors"
	"time"
)

var ErrNewQuestionBeforeAnswer = errors.New("previous question hasn't been answered yet")

type Conversation struct {
	QAs         []qa      `json:"qa,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProfileSlug string    `json:"profile_slug"`
	ModelSlug   string    `json:"model_slug"`
}

func NewConversation(profileSlug string, modelSlug string) *Conversation {
	creation := time.Now()
	return &Conversation{
		QAs:         make([]qa, 0),
		CreatedAt:   creation,
		UpdatedAt:   creation,
		ModelSlug:   modelSlug,
		ProfileSlug: profileSlug,
	}
}

func (c *Conversation) Clear() {
	c.QAs = make([]qa, 0)
	c.markAsUpdated()
}

func (c *Conversation) Len() int {
	return len(c.QAs)
}

func (c *Conversation) LastAnswer() string {
	if len(c.QAs) == 0 {
		return ""
	}
	return c.QAs[len(c.QAs)-1].GetAnswer()
}

func (c *Conversation) LastQuestion() string {
	if len(c.QAs) == 0 {
		return ""
	}
	return c.QAs[len(c.QAs)-1].GetQuestion()
}

func (c *Conversation) NewQuestion(question string) error {
	if len(c.QAs) > 0 && c.LastAnswer() == "" {
		return ErrNewQuestionBeforeAnswer
	}

	c.QAs = append(c.QAs, *NewQA(question))
	c.markAsUpdated()
	return nil
}

func (c *Conversation) AnswerLastQuestion(answer string) error {
	c.markAsUpdated()
	return c.QAs[len(c.QAs)-1].addAnswer(answer)
}

func (c *Conversation) SetProfileTo(profileSlug string) {
	c.ProfileSlug = profileSlug
}

func (c *Conversation) SetModelTo(modelSlug string) {
	c.ModelSlug = modelSlug
}

func (c *Conversation) markAsUpdated() {
	c.UpdatedAt = time.Now()
}
