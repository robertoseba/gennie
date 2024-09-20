package chat

type ChatHistory struct {
	Responses   []Response
	persistence IPersistence
}

func (c ChatHistory) lastResponse() Response {
	if len(c.Responses) == 0 {
		return Response{}
	}
	return c.Responses[len(c.Responses)-1]
}

func (c *ChatHistory) addResponse(response Response) {
	c.Responses = append(c.Responses, response)
}

func (c *ChatHistory) saveToDisk() error {
	return c.persistence.save(*c)
}

func (c *ChatHistory) loadFromDisk() error {
	history, err := c.persistence.load()
	if err != nil {
		return err
	}
	c.Responses = history.Responses
	return nil
}

func NewChatHistory(name string, persistence IPersistence) *ChatHistory {
	return &ChatHistory{
		persistence: persistence,
	}
}
