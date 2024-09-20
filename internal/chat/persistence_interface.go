package chat

// Interface for persisting the chat history
type IPersistence interface {
	save(history ChatHistory) error
	load() (ChatHistory, error)
}
