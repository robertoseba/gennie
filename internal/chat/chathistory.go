package chat

type ChatHistory struct {
	Responses []Chat
}

func NewChatHistory() ChatHistory {
	return ChatHistory{
		Responses: nil,
	}
}

// Returns the last chat in the chat history. If there are no responses, it returns false with empty Chat.
func (c ChatHistory) LastChat() (Chat, bool) {
	if len(c.Responses) == 0 {
		return Chat{}, false
	}
	return c.Responses[len(c.Responses)-1], true
}

/**
 * Adds a chat to the chat history.
 * The chat can can still be incomplete, but it must have at least a question.
 * The answer will be added later by the model.
 */
func (c *ChatHistory) AddChat(chat Chat) {
	c.Responses = append(c.Responses, chat)
}

func (c *ChatHistory) Clear() {
	c.Responses = nil
}

func (c *ChatHistory) Len() int {
	return len(c.Responses)
}

func (c *ChatHistory) LastAnswer() string {
	if len(c.Responses) == 0 {
		return ""
	}
	return c.Responses[len(c.Responses)-1].GetAnswer()
}

func (c *ChatHistory) LastQuestion() string {
	if len(c.Responses) == 0 {
		return ""
	}
	return c.Responses[len(c.Responses)-1].GetQuestion()
}

func (c *ChatHistory) SetNewAnswerToLastChat(answer string) {
	c.Responses[len(c.Responses)-1].AddAnswer(answer)
}
