package cmd

import (
	"slices"
	"testing"
)

func TestHasConfigCommands(t *testing.T) {
	r := NewConfigCmd(nil, nil)
	if r.Use != "config" {
		t.Errorf("Expected 'config' but got %s", r.Use)
	}
	expectedCommands := []string{"show", "model"}

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
