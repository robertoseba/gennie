package cache

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/robertoseba/gennie/internal/conversation"
)

func TestStorageDefaultPath(t *testing.T) {

	storePath, err := GetStorageFilepath()
	expectedPath := "/gennie/.cache"

	if err != nil {
		t.Error("Error while getting cache path")
	}
	if !strings.HasSuffix(storePath, expectedPath) {
		t.Errorf("Default storage path should end with: /gennie/.cache, got: %s", storePath)
	}
}

func TestRestoreFromNonExistentFile(t *testing.T) {
	_, err := RestoreFrom("non_existent_file")
	if !errors.Is(err, ErrNoStoreFile) {
		t.Errorf("Expected error to be ErrNoCacheFile, got: %v", err)
	}
}

func TestRestoreFrom(t *testing.T) {
	cachePath := ".cache_temp"
	c := NewStorage(cachePath)

	c.ChatHistory = conversation.NewConversation()
	c.ChatHistory.NewQuestion("question")

	c.CurrModelSlug = "testModelSlug"

	c.Save()

	restoredCache, err := RestoreFrom(cachePath)

	if err != nil {
		t.Errorf("Error while restoring cache: %v", err)
	}

	if restoredCache.filePath != cachePath {
		t.Errorf("Expected restored cache file path to be %s, got: %s", cachePath, restoredCache.filePath)
	}

	if restoredCache.Config != c.Config {
		t.Errorf("Expected restored cache Config to be %v, got: %v", c.Config, restoredCache.Config)
	}

	if restoredCache.CachedProfiles == nil {
		t.Errorf("Expected restored cache CachedProfilesPath to be %v, got: %v", c.CachedProfiles, restoredCache.CachedProfiles)
	}

	if restoredCache.ChatHistory.LastQuestion() != "question" {
		t.Errorf("Expected restored cache ChatHistory to be %v, got: %v", c.ChatHistory, restoredCache.ChatHistory)
	}

	if restoredCache.IsNew() {
		t.Errorf("Expected restored cache to not be new")
	}

	os.Remove(cachePath)

}
