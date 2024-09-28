package cmd

import (
	"testing"
)

func TestHasAllFlags(t *testing.T) {
	r := NewAskCmd(nil, nil, nil)
	if r.Use != "ask [question for the llm model]" {
		t.Errorf("Expected 'ask' but got %s", r.Use)
	}
	expectedFlags := []string{"followup", "append", "model", "profile"}

	for _, f := range expectedFlags {
		if r.Flags().Lookup(f) == nil {
			t.Errorf("Expected flag %s not found", f)
		}
	}

}
