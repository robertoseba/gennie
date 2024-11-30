package profile

import (
	"testing"

	"github.com/robertoseba/gennie/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestListAll(t *testing.T) {
	config := config.NewConfig()
	pr := NewProfileRepository(&config, map[string]ProfileInfo{
		"test": {
			Name:     "Test Profile",
			Slug:     "test",
			Filepath: "./test.profile.toml",
		},
		"test2": {
			Name:     "Test Profile 2",
			Slug:     "test2",
			Filepath: "./test2.profile.toml",
		},
	})

	t.Run("List all profiles", func(t *testing.T) {
		allProfiles := pr.ListAll()

		assert.Equal(t, 2, len(allProfiles))
		assert.Equal(t, "Test Profile", allProfiles["test"])
		assert.Equal(t, "Test Profile 2", allProfiles["test2"])
	})

	t.Run("List all profiles with no profiles", func(t *testing.T) {
		pr = NewProfileRepository(&config, map[string]ProfileInfo{})
		allProfiles := pr.ListAll()

		assert.Equal(t, 0, len(allProfiles))
	})
}

func TestFindBySlug(t *testing.T) {
	config := config.NewConfig()
	pr := NewProfileRepository(&config, map[string]ProfileInfo{
		"test": {
			Name:     "Test Profile",
			Slug:     "test",
			Filepath: "./stub/stub.profile.toml",
		},
		"test2": {
			Name:     "Test Profile 2",
			Slug:     "test2",
			Filepath: "./test2.profile.toml",
		},
	})

	t.Run("Find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("test")

		assert.Nil(t, err)
		assert.Equal(t, "Test Profile", profile.Name)
		assert.Equal(t, "test", profile.Slug)
		assert.Equal(t, "Roberto Seba", profile.Author)
		assert.Equal(t, "just a profile stub for testing", profile.Data)
	})

	t.Run("Cant find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("test3")
		assert.EqualError(t, err, "profile not found. Try using refresh command if you're sure the profile exists")
		assert.Nil(t, profile)
	})
}

func TestRefreshProfiles(t *testing.T) {
	t.Run("Refresh profiles", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir("./stub")

		pr := NewProfileRepository(&config, map[string]ProfileInfo{})
		err := pr.RefreshProfiles()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(pr.profilesInfo))

		assert.Equal(t, "Test Profile", pr.profilesInfo["test"].Name)
	})

	t.Run("Fails when dir doest not exist", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir("./stub/invalid")

		pr := NewProfileRepository(&config, map[string]ProfileInfo{})
		err := pr.RefreshProfiles()

		assert.EqualError(t, err, "error reading profile directory")
	})

	t.Run("empty when no profiles found", func(t *testing.T) {
		config := config.NewConfig()
		config.SetProfilesDir(".")

		pr := NewProfileRepository(&config, map[string]ProfileInfo{})
		err := pr.RefreshProfiles()

		assert.Nil(t, err)
		assert.Equal(t, 0, len(pr.profilesInfo))
	})
}
