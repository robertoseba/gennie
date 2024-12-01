package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAll(t *testing.T) {

	t.Run("List all profiles including default", func(t *testing.T) {
		pr := NewProfileRepository("./stub")
		allProfiles, err := pr.ListAll()

		assert.Nil(t, err)
		assert.Equal(t, 3, len(allProfiles))
		assert.Equal(t, "Test Profile", allProfiles["stub"].Name)
		assert.Equal(t, "Test Profile 2", allProfiles["stub2"].Name)
		assert.Equal(t, "Default assistant", allProfiles["default"].Name)
	})

	t.Run("when no profiles available lists only default", func(t *testing.T) {
		pr := NewProfileRepository(".")
		allProfiles, err := pr.ListAll()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(allProfiles))
		assert.Equal(t, "Default assistant", allProfiles["default"].Name)
	})
	t.Run("fails if cant find profile dir", func(t *testing.T) {
		pr := NewProfileRepository("./invalid")
		allProfiles, err := pr.ListAll()

		assert.EqualError(t, err, "error reading profiles directory: open ./invalid: no such file or directory")
		assert.Nil(t, allProfiles)
	})
}

func TestFindBySlug(t *testing.T) {

	pr := NewProfileRepository("./stub")
	t.Run("Find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("stub2")

		assert.Nil(t, err)
		assert.Equal(t, "Test Profile 2", profile.Name)
		assert.Equal(t, "stub2", profile.Slug)
		assert.Equal(t, "Roberto Seba", profile.Author)
		assert.Equal(t, "just a profile stub for testing - number 2", profile.Data)
	})

	t.Run("Cant find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("test3")
		assert.EqualError(t, err, "error loading toml file: open stub/test3.profile.toml: no such file or directory")
		assert.Nil(t, profile)
	})
}
