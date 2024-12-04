package usecases

import (
	"testing"
)

func TestGetAnswerService(t *testing.T) {
	//TODO: create mock IAPIclient
	t.Run("completes the conversation with answers from the API", func(t *testing.T) {
	})

	t.Run("returns an error if cant find profile", func(t *testing.T) {
	})

	t.Run("returns an error if cant find model", func(t *testing.T) {
	})

	t.Run("appends the content of a file to the conversation", func(t *testing.T) {
	})

	t.Run("returns an error if the conversation cannot be loaded", func(t *testing.T) {
	})

	t.Run("when model is inputed replaces the model in active conversation", func(t *testing.T) {
	})

	t.Run("when profile is inputed replaces the profile in active conversation", func(t *testing.T) {
	})

	t.Run("when input is a not set as follow up question, creates a new conversation", func(t *testing.T) {
	})

	t.Run("when input is a follow up question, appends the question to the conversation", func(t *testing.T) {
	})
}
