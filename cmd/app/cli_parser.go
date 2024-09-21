package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/robertoseba/gennie/internal/models"
)

type InputOptions struct {
	Model    models.ModelEnum
	Question string
}

func ParseCliOptions() *InputOptions {

	question := strings.Join(os.Args[1:], " ")

	if question == "" {
		wrongUsage()
	}

	return &InputOptions{
		Model:    models.OpenAI,
		Question: question,
	}
}

func wrongUsage() {
	fmt.Println("Please provide a question to ask.")
	os.Exit(1)
}
