package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

var ErrNoActiveConversation = errors.New("no active conversation")
var ErrConversationNotFound = errors.New("conversation not found")

const ActiveConversationFileName = "active.json"

type ConversationRepository struct {
	cacheDir string
}

func NewConversationRepository(cacheDir string) *ConversationRepository {
	return &ConversationRepository{cacheDir: cacheDir}
}

// Loads the last conversation that has been active
// If there is no active conversation, it returns an error
func (r *ConversationRepository) LoadActive() (*conversation.Conversation, error) {
	c, err := r.loadFrom(path.Join(r.cacheDir, ActiveConversationFileName))

	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			return nil, ErrNoActiveConversation
		}
		return nil, err
	}

	return c, nil
}

func (r *ConversationRepository) LoadFromFile(filepath string) (*conversation.Conversation, error) {
	return r.loadFrom(filepath)
}

func (r *ConversationRepository) ExportToFile(conversation *conversation.Conversation, filepath string) error {
	content, err := json.Marshal(conversation)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, content, 0644)
}

func (r *ConversationRepository) SaveAsActive(conversation *conversation.Conversation) error {
	filepath := path.Join(r.cacheDir, ActiveConversationFileName)

	return r.ExportToFile(conversation, filepath)
}

func (r *ConversationRepository) loadFrom(filepath string) (*conversation.Conversation, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, ErrConversationNotFound
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	conversation := &conversation.Conversation{}

	err = json.Unmarshal(content, conversation)

	if err != nil {
		return nil, fmt.Errorf("error decoding conversation: %w", err)
	}

	return conversation, nil
}

//TODO: CreateCacheDir func to create the cache dir if it does not exist using os.UserCacheDir
