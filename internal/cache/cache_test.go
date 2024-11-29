package cache

import (
	"errors"
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/profile"
)

func TestNewStorageHasNewConfig(t *testing.T) {
	c := NewStorage(".cache")
	expectedConfig := common.NewConfig()

	if c.Config != expectedConfig {
		t.Errorf("Expected new storage to have new config, got: %v", c.Config)
	}
}

func TestNewStorageSetAsNew(t *testing.T) {
	c := NewStorage(".cache")

	if !c.isNew {
		t.Errorf("Storage should have been set as new")
	}
}

func TestNewStorageHasFilepathSet(t *testing.T) {
	c := NewStorage("/temp/.cache")

	if c.filePath != "/temp/.cache" {
		t.Errorf("Expected cache filepath to be .cache, got: %s", c.filePath)
	}
}

func TestSavesStorageToFile(t *testing.T) {
	c := NewStorage(".cache_temp")

	err := c.Save()
	if err != nil {
		t.Errorf("Error while saving cache: %v", err)
	}

	if _, err := os.Stat(".cache_temp"); os.IsNotExist(err) {
		t.Errorf("Expected cache file to exist, got: %v", err)
	}

	os.Remove(".cache_temp")
}

func TestGetProfileSlugs(t *testing.T) {
	c := NewStorage(".cache_temp")

	c.CachedProfiles = map[string]profile.ProfileInfo{
		"test": {
			Slug:     "test",
			Name:     "test",
			Filepath: "/test.profile.toml",
		},
		"test2": {
			Slug:     "test2",
			Name:     "test2",
			Filepath: "/test2.profile.toml",
		},
	}

	profileSlugs := c.GetCachedProfiles()

	if len(profileSlugs) != 2 {
		t.Errorf("Expected 2 profile slugs, got: %v", profileSlugs)
	}

	for slug, pInfo := range profileSlugs {
		if slug != "test" && slug != "test2" {
			t.Errorf("Expected profile slug to be test or test2, got: %s", slug)
		}
		if pInfo.Filepath != "/test.profile.toml" && pInfo.Filepath != "/test2.profile.toml" {
			t.Errorf("Expected profile filepath to be /test.profile.toml or /test2.profile.toml, got: %s", pInfo.Filepath)
		}
		if pInfo.Name != "test" && pInfo.Name != "test2" {
			t.Errorf("Expected profile name to be test or test2, got: %s", pInfo.Name)
		}
	}
}

func TestLoadsProfileDataFromFile(t *testing.T) {
	c := NewStorage(".cache_temp")

	c.CachedProfiles = map[string]profile.ProfileInfo{
		"test": {
			Slug:     "test",
			Name:     "test",
			Filepath: "/test.profile.toml",
		},
		"stub": {
			Slug:     "stub",
			Name:     "profileStub",
			Filepath: "./test/stub.profile.toml",
		},
	}

	profile, err := c.LoadProfileData("stub")
	if err != nil {
		t.Errorf("Error while getting profile: %v", err)
	}

	if profile.Name != "profileStub" {
		t.Errorf("Expected profile name be %s, got: %s", "profileStub", profile.Name)
	}

	if profile.Data != "just a profile stub for testing" {
		t.Errorf("Got wrong profile data then expected: %s", profile.Data)
	}
}

func TestGetProfileWithInexistentSlug(t *testing.T) {
	c := NewStorage(".cache_temp")

	_, err := c.LoadProfileData("inexistent")

	if !errors.Is(err, ErrNoProfileSlug) {
		t.Errorf("Expected ErrNoProfileSlug, got: %v", err)
	}
}

func TestGetProfileForDefaultSlug(t *testing.T) {
	c := NewStorage(".cache_temp")

	profileRetrieved, err := c.LoadProfileData("default")
	if err != nil {
		t.Errorf("Error while getting profile: %v", err)
	}

	defaultExpected := profile.DefaultProfile()

	if *profileRetrieved != *defaultExpected {
		t.Errorf("Expected default profile, got: %v", profileRetrieved)
	}
}

func TestClear(t *testing.T) {
	c := NewStorage(".cache_temp")

	c.CurrModelSlug = "newSlug"

	c.CachedProfiles = map[string]profile.ProfileInfo{
		"test": {
			Slug:     "test",
			Name:     "test",
			Filepath: "/test.profile.toml",
		},
		"test2": {
			Slug:     "test2",
			Name:     "test2",
			Filepath: "/test2.profile.toml",
		},
	}

	chat := chat.NewQA("testing question")
	chat.AddAnswer("testing answer")
	c.ChatHistory.AddChat(*chat)

	c.Config.AnthropicApiKey = "test_key"

	c.Clear()

	if len(c.CachedProfiles) != 0 {
		t.Errorf("Expected no cached profiles, got: %v", c.CachedProfiles)
	}

	if c.ChatHistory.Len() != 0 {
		t.Errorf("Expected no chat history, got: %v", c.ChatHistory)
	}

	if c.Config.AnthropicApiKey != "" {
		t.Errorf("Expected no api key, got: %s", c.Config.AnthropicApiKey)
	}

	if c.CurrProfile != *profile.DefaultProfile() {
		t.Errorf("Expected default profile, got: %v", c.CurrProfile)
	}

	if c.CurrModelSlug != "default" {
		t.Errorf("Expected default model slug, got: %s", c.CurrModelSlug)
	}
}
