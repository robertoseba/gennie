package cmd

import (
	"slices"
	"testing"
)

func TestHasProfileSubCommands(t *testing.T) {
	r := NewProfilesCmd(nil, nil)
	if r.Use != "profile" {
		t.Errorf("Expected 'profiles' but got %s", r.Use)
	}
	expectedCommands := []string{"list", "refresh", "config"}

	for _, c := range r.Commands() {
		idx := slices.Index(expectedCommands, c.Use)
		if idx == -1 {
			t.Errorf("Expected command %s not found", c.Use)
			continue
		}

		expectedCommands = slices.Delete(expectedCommands, idx, idx+1)
	}

	if len(expectedCommands) > 0 {
		t.Errorf("Missing commands: %v", expectedCommands)
	}
}
