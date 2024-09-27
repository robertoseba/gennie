package chat

type ChatHistory struct {
	Responses []Chat
}

// Returns the last response in the chat history. If there are no responses, it returns false.
func (c ChatHistory) LastResponse() (Chat, bool) {
	if len(c.Responses) == 0 {
		return Chat{}, false
	}
	return c.Responses[len(c.Responses)-1], true
}

func (c *ChatHistory) AddResponse(response Chat) {
	// if c.Responses == nil {
	// 	c.Responses = make([]Chat, 0, 1)
	// }
	c.Responses = append(c.Responses, response)
}

func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		Responses: nil,
	}
}

func (c *ChatHistory) Clear() {
	c.Responses = nil
}
