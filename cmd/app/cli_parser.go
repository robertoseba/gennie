package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/robertoseba/gennie/internal/models"
)

type InputOptions struct {
	Model      models.ModelEnum
	Question   string
	ConfigMode bool
}

func ParseCliOptions() *InputOptions {
	config := flag.Bool("config", false, "Activate the configuration mode")

	if len(os.Args) <= 1 {
		wrongUsage()
	}

	sortCliOptions()

	flag.Parse()

	fmt.Println(flag.Args())
	question := strings.Join(os.Args[1:], " ")
	fmt.Println(question)
	if question == "" {
		wrongUsage()
	}
	os.Exit(0)
	return &InputOptions{
		Model:      models.OpenAI,
		Question:   question,
		ConfigMode: *config,
	}
}

func wrongUsage() {
	fmt.Println("Please provide a question to ask.")
	os.Exit(1)
}

func sortCliOptions() {
	if os.Args[1][0] == '-' {
		// In case the user passes the filename as the first argument before flags, we need to move it to the end
		// so that the flag package can parse the flags correctly
		question := os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		os.Args = append(os.Args, question)
	}
}
