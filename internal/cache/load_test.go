package cache

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/robertoseba/gennie/internal/chat"
)

func TestCacheDefaultPath(t *testing.T) {

	cachePath, err := GetCacheFilePath()
	expectedPath := "/gennie/.cache"

	if err != nil {
		t.Error("Error while getting cache path")
	}
	if !strings.HasSuffix(cachePath, expectedPath) {
		t.Errorf("Default cache path should end with: /gennie/.cache, got: %s", cachePath)
	}
}

func TestRestoreFromNonExistentFile(t *testing.T) {
	_, err := RestoreFrom("non_existent_file")
	if !errors.Is(err, ErrNoCacheFile) {
		t.Errorf("Expected error to be ErrNoCacheFile, got: %v", err)
	}
}

func TestRestoreFrom(t *testing.T) {
	cachePath := ".cache_temp"
	c := NewCache(cachePath)

	chatHistory := chat.NewChatHistory()
	chat := chat.NewChat("question")
	chatHistory.AddChat(*chat)
	c.ChatHistory = chatHistory

	c.Config.CurrModelSlug = "testModelSlug"

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

	if restoredCache.CachedProfilesPath == nil {
		t.Errorf("Expected restored cache CachedProfilesPath to be %v, got: %v", c.CachedProfilesPath, restoredCache.CachedProfilesPath)
	}

	if restoredCache.ChatHistory.LastQuestion() != "question" {
		t.Errorf("Expected restored cache ChatHistory to be %v, got: %v", c.ChatHistory, restoredCache.ChatHistory)
	}

	os.Remove(cachePath)

}
