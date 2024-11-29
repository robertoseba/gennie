package conversation

import "errors"

type Conversation struct {
	QAs []qa
}

func NewConversation() Conversation {
	return Conversation{
		QAs: make([]qa, 0),
	}
}

/**
 * Returns the last question/answer in the conversation.
 * If there are no responses, it returns false with empty QA.
 */
//TODO: remove lastQA and use LastQuestion and LastAnswer
func (c Conversation) LastQA() (qa, bool) {
	if len(c.QAs) == 0 {
		return qa{}, false
	}
	return c.QAs[len(c.QAs)-1], true
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
		return errors.New("previous question hasn't been answered yet")
	}
	c.QAs = append(c.QAs, *NewQA(question))
	return nil
}

func (c *Conversation) AnswerLastQuestion(answer string) error {
	return c.QAs[len(c.QAs)-1].AddAnswer(answer)
}

//TODO: create encoder/decoder for conversation
