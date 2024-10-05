package cache

import (
	"errors"
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/common"
)

func TestNewCacheHasNewConfig(t *testing.T) {
	c := NewCache(".cache")
	expectedConfig := common.NewConfig()

	if c.Config != expectedConfig {
		t.Errorf("Expected new cache to have new config, got: %v", c.Config)
	}
}

func TestNewCacheHasFilepathSet(t *testing.T) {
	c := NewCache("/temp/.cache")

	if c.filePath != "/temp/.cache" {
		t.Errorf("Expected cache filepath to be .cache, got: %s", c.filePath)
	}
}

func TestSavesCacheToFile(t *testing.T) {
	c := NewCache(".cache_temp")

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
	c := NewCache(".cache_temp")

	c.CachedProfilesPath = map[string]string{
		"test":  "/test.profile.json",
		"test2": "/test2.profile.json",
	}

	profileSlugs := c.GetProfileSlugs()

	if len(profileSlugs) != 2 {
		t.Errorf("Expected 2 profile slugs, got: %v", profileSlugs)
	}

	for _, slug := range profileSlugs {
		if slug != "test" && slug != "test2" {
			t.Errorf("Expected profile slug to be test or test2, got: %s", slug)
		}
	}
}

func TestGetProfile(t *testing.T) {
	c := NewCache(".cache_temp")

	c.CachedProfilesPath = map[string]string{
		"stub": "./stub/stub.profile.json",
	}

	profile, err := c.GetProfile("stub")

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
	c := NewCache(".cache_temp")

	_, err := c.GetProfile("inexistent")

	if !errors.Is(err, ErrNoProfileSlug) {
		t.Errorf("Expected ErrNoProfileSlug, got: %v", err)
	}
}
