package conversation

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
func (c Conversation) LastQA() (qa, bool) {
	if len(c.QAs) == 0 {
		return qa{}, false
	}
	return c.QAs[len(c.QAs)-1], true
}

/**
 * Adds a question/answer to the conversation.
 * The QAs can can still be incomplete, but it must have at least a question.
 * The answer will be added later by the model.
 */
func (c *Conversation) AddQA(qa qa) {
	c.QAs = append(c.QAs, qa)
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

func (c *Conversation) SetNewAnswerToLastChat(answer string) error {
	return c.QAs[len(c.QAs)-1].AddAnswer(answer)
}
