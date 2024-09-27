package chat

type ChatHistory struct {
	Responses []*Chat
}

func (c ChatHistory) LastResponse() *Chat {
	if len(c.Responses) == 0 {
		return nil
	}
	return c.Responses[len(c.Responses)-1]
}

func (c *ChatHistory) AddResponse(response *Chat) {
	c.Responses = append(c.Responses, response)
}

func NewChatHistory(name string) *ChatHistory {
	return &ChatHistory{}
}
