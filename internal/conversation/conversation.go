package conversation

import "errors"

var ErrNewQuestionBeforeAnswer = errors.New("previous question hasn't been answered yet")

type Conversation struct {
	QAs []qa
}

func NewConversation() Conversation {
	return Conversation{
		QAs: make([]qa, 0),
	}
}

func (c *Conversation) Clear() {
	c.QAs = make([]qa, 0)
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
	return nil
}

func (c *Conversation) AnswerLastQuestion(answer string) error {
	return c.QAs[len(c.QAs)-1].addAnswer(answer)
}

//TODO: create encoder/decoder for conversation so we dont have to expose QA?
