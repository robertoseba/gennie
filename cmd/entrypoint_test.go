package cmd

import (
	"slices"
	"testing"
)

func TestHasRootSubCommands(t *testing.T) {
	r := NewRootCmd(nil, nil, nil)
	if r.Use != "gennie" {
		t.Errorf("Expected 'gennie' but got %s", r.Use)
	}
	expectedCommands := []string{"model", "profile", "ask [question for the llm model]", "status", "export", "clear"}

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
