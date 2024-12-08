package cmd

import (
	"bytes"
	"slices"
	"testing"

	"github.com/robertoseba/gennie/internal/infra/container"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestHasRootSubCommands(t *testing.T) {
	t.Run("keyword is gennie", func(t *testing.T) {
		r, _, _ := setupCommand(t)
		require.Equal(t, "gennie", r.Use)
	})

	t.Run("decriptions", func(t *testing.T) {
		r, stdout, _ := setupCommand(t)
		r.SetArgs([]string{"--help"})
		r.Execute()
		require.Contains(t, stdout.String(), "Gennie is a cli assistant with multiple models and profile support.")
	})

	t.Run("template version", func(t *testing.T) {
		r, stdout, _ := setupCommand(t)
		r.SetArgs([]string{"--version"})
		r.Execute()
		require.Equal(t, "Gennie version: 0.0.1", stdout.String())
	})

	t.Run("sub commands are", func(t *testing.T) {
		r, _, _ := setupCommand(t)
		expectedCommands := []string{"config", "model", "profile", "ask [question for the llm model]", "status", "conversation [command]"}

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
	})
}

func setupCommand(t *testing.T) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Setenv("XDG_CONFIG_HOME", ".tmp")
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	r := newRootCmd("0.0.1", stdOut, stdErr)
	setupSubCommands(r, container.NewContainer(), nil)
	return r, stdOut, stdErr
}
