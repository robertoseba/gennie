package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type InputOptions struct {
	Question   string
	ConfigMode bool
}

func ParseCliOptions() *InputOptions {
	config := flag.Bool("config", false, "Activates configuration mode")

	if len(os.Args) <= 1 {
		wrongUsage()
	}

	var question string

	if isQuestion() {
		question = strings.Join(os.Args[1:], " ")

		if question == "" {
			wrongUsage()
		}
	} else {
		flag.Parse()
	}

	return &InputOptions{
		Question:   question,
		ConfigMode: *config,
	}
}

func wrongUsage() {
	fmt.Println("Please provide a question to ask or use one of the options below:")
	flag.PrintDefaults()
	os.Exit(1)
}

func isQuestion() bool {
	firstChars := os.Args[1][0:2]
	if firstChars != "--" {
		return true
	}
	return false
}
