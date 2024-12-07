package repositories

import (
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/core/profile"
	"github.com/stretchr/testify/require"
)

func TestListAll(t *testing.T) {

	t.Run("List all profiles including default", func(t *testing.T) {
		pr := NewProfileRepository("./stub")
		allProfiles, err := pr.ListAll()

		require.NoError(t, err)
		require.Len(t, allProfiles, 3)
		require.Equal(t, "Test Profile", allProfiles["stub"].Name)
		require.Equal(t, "Test Profile 2", allProfiles["stub2"].Name)
		require.Equal(t, "Default assistant", allProfiles["default"].Name)
	})

	t.Run("when no profiles available lists only default", func(t *testing.T) {
		pr := NewProfileRepository(".")
		allProfiles, err := pr.ListAll()

		require.NoError(t, err)
		require.Len(t, allProfiles, 1)
		require.Equal(t, "Default assistant", allProfiles["default"].Name)
	})
	t.Run("fails if cant find profile dir", func(t *testing.T) {
		pr := NewProfileRepository("./invalid")
		allProfiles, err := pr.ListAll()

		require.EqualError(t, err, "no profiles found. Please add profiles to the profiles folder.")
		require.Len(t, allProfiles, 1)
		require.Equal(t, profile.DefaultProfile().Name, allProfiles[profile.DefaultProfileSlug].Name)
	})
}

func TestFindBySlug(t *testing.T) {

	pr := NewProfileRepository("./stub")
	t.Run("Find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("stub2")

		require.NoError(t, err)
		require.Equal(t, "Test Profile 2", profile.Name)
		require.Equal(t, "stub2", profile.Slug)
		require.Equal(t, "Roberto Seba", profile.Author)
		require.Equal(t, "just a profile stub for testing - number 2", profile.Data)
	})

	t.Run("Cant find by slug", func(t *testing.T) {
		profile, err := pr.FindBySlug("test3")
		require.EqualError(t, err, "error loading toml file: open stub/test3.profile.toml: no such file or directory")
		require.Nil(t, profile)
	})
}

func TestDefaultProfileDir(t *testing.T) {
	t.Run("Default profile dir", func(t *testing.T) {
		p := DefaultProfileDir()
		osConfigDir, err := os.UserConfigDir()
		require.NoError(t, err)

		require.Equal(t, osConfigDir+"/gennie/profiles", p)
	})
}
