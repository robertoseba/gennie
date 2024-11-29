package chat

type ChatHistory struct {
	QAs []QA
}

func NewChatHistory() ChatHistory {
	return ChatHistory{
		QAs: make([]QA, 0),
	}
}

/**
 * Returns the last question/answer in the history.
 * If there are no responses, it returns false with empty QA.
 */
func (c ChatHistory) LastQA() (QA, bool) {
	if len(c.QAs) == 0 {
		return QA{}, false
	}
	return c.QAs[len(c.QAs)-1], true
}

/**
 * Adds a question/answer to the history.
 * The QAs can can still be incomplete, but it must have at least a question.
 * The answer will be added later by the model.
 */
func (c *ChatHistory) AddQA(qa QA) {
	c.QAs = append(c.QAs, qa)
}

func (c *ChatHistory) Clear() {
	c.QAs = make([]QA, 0)
}

func (c *ChatHistory) Len() int {
	return len(c.QAs)
}

func (c *ChatHistory) LastAnswer() string {
	if len(c.QAs) == 0 {
		return ""
	}
	return c.QAs[len(c.QAs)-1].GetAnswer()
}

func (c *ChatHistory) LastQuestion() string {
	if len(c.QAs) == 0 {
		return ""
	}
	return c.QAs[len(c.QAs)-1].GetQuestion()
}

func (c *ChatHistory) SetNewAnswerToLastChat(answer string) error {
	return c.QAs[len(c.QAs)-1].AddAnswer(answer)
}
