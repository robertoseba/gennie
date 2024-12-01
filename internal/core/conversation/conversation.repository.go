package conversation

type ConversationRepository interface {
	LoadActive() (*Conversation, error)
	LoadFromFile(filepath string) (*Conversation, error)
	ExportToFile(conversation *Conversation, filepath string) error
	SaveAsActive(conversation *Conversation) error
}
