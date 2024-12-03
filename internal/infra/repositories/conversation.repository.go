package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/profile"
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
// If there is no active conversation, creates a new one with the default profile and model
func (r *ConversationRepository) LoadActive() (*conversation.Conversation, error) {
	//TODO: cache conversation loaded
	c, err := r.loadFrom(path.Join(r.cacheDir, ActiveConversationFileName))

	if err != nil {
		if errors.Is(err, ErrConversationNotFound) {
			return conversation.NewConversation(profile.DefaultProfileSlug, string(models.DefaultModel)), nil
		}
		return nil, err
	}

	return c, nil
}

func (r *ConversationRepository) LoadFromFile(filepath string) (*conversation.Conversation, error) {
	return r.loadFrom(filepath)
}

func (r *ConversationRepository) ExportToFile(conversation *conversation.Conversation, filepath string) error {
	content, err := json.MarshalIndent(conversation, "", "  ")
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
