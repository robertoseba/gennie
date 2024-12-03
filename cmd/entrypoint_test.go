package cmd

import (
	"testing"
)

func TestHasRootSubCommands(t *testing.T) {
	t.Skip("Skipping test for now")
	// cmdUtil, _ := NewCmdUtil("0.0.1")
	// r := NewRootCmd(cmdUtil)

	// if r.Use != "gennie" {
	// 	t.Errorf("Expected 'gennie' but got %s", r.Use)
	// }
	// expectedCommands := []string{"config", "model", "profile", "ask [question for the llm model]", "status", "export [filename]", "clear"}

	// for _, c := range r.Commands() {
	// 	idx := slices.Index(expectedCommands, c.Use)
	// 	if idx == -1 {
	// 		t.Errorf("Expected command %s not found", c.Use)
	// 		continue
	// 	}

	// 	expectedCommands = slices.Delete(expectedCommands, idx, idx+1)
	// }

	// if len(expectedCommands) > 0 {
	// 	t.Errorf("Missing commands: %v", expectedCommands)
	// }
}
