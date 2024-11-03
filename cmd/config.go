package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/profile"
	"github.com/spf13/cobra"
)

func NewConfigCmd(storage common.IStorage, p *output.Printer) *cobra.Command {

	clearCmd := &cobra.Command{
		Use:   "config",
		Short: "Configures Gennie",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := storage.GetConfig()

			configApiKeys(p, &config)
			configProfile(p, &config)
			configOllama(p, &config)

			storage.SetConfig(config)

			refreshProfiles(storage)

			return nil
		},
	}

	return clearCmd
}

func configApiKeys(p *output.Printer, config *common.Config) {
	config.OpenAiApiKey = askKey(p, "OpenAI", config.OpenAiApiKey)
	config.AnthropicApiKey = askKey(p, "Anthropic", config.AnthropicApiKey)
	config.MaritacaApiKey = askKey(p, "Maritaca", config.MaritacaApiKey)
	config.GroqApiKey = askKey(p, "Groq", config.GroqApiKey)
}

func askKey(p *output.Printer, key string, previousValue string) string {
	question := output.NewQuestion(fmt.Sprintf("Enter your %s API Key", key))

	if previousValue != "" {
		question.WithPrevious(previousValue, true)
	}
	return question.Ask(p)
}

func configProfile(p *output.Printer, config *common.Config) {
	previousValue := config.ProfilesPath
	question := output.NewQuestion("Enter your profiles folder path")

	if previousValue != "" {
		question.WithPrevious(previousValue, false)
	} else {
		question.WithPrevious(profile.DefaultProfilePath(), false)
	}
	config.ProfilesPath = question.Ask(p)
}

func configOllama(p *output.Printer, config *common.Config) {
	q := output.NewQuestion("What is your Ollama host address?")
	if config.OllamaHost != "" {
		q.WithPrevious(config.OllamaHost, false)
	}
	config.OllamaHost = q.Ask(p)

	q = output.NewQuestion("What Ollama model would you like to use?")
	if config.OllamaModel != "" {
		q.WithPrevious(config.OllamaModel, false)
	}
	config.OllamaModel = q.Ask(p)
}
