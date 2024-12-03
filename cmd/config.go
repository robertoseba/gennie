package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/infra/repositories"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewConfigCmd(configRepo config.IConfigRepository, p *output.Printer) *cobra.Command {

	clearCmd := &cobra.Command{
		Use:   "config",
		Short: "Configures Gennie",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := configRepo.Load()
			if err != nil {
				return err
			}

			configApiKeys(p, config)
			configProfile(p, config)
			configOllama(p, config)

			return configRepo.Save(config)
		},
	}

	return clearCmd
}

func configApiKeys(p *output.Printer, config *config.Config) {
	config.APIKeys.OpenAiApiKey = askKey(p, "OpenAI", config.APIKeys.OpenAiApiKey)
	config.APIKeys.AnthropicApiKey = askKey(p, "Anthropic", config.APIKeys.AnthropicApiKey)
	config.APIKeys.MaritacaApiKey = askKey(p, "Maritaca", config.APIKeys.MaritacaApiKey)
	config.APIKeys.GroqApiKey = askKey(p, "Groq", config.APIKeys.GroqApiKey)
}

func askKey(p *output.Printer, key string, previousValue string) string {
	question := output.NewQuestion(fmt.Sprintf("Enter your %s API Key", key))

	if previousValue != "" {
		question.WithPrevious(previousValue, true)
	}
	return question.Ask(p)
}

func configProfile(p *output.Printer, config *config.Config) {
	question := output.NewQuestion("Enter your profiles folder path")

	previousValue := config.ProfilesDirPath

	if previousValue != "" {
		question.WithPrevious(previousValue, false)
	} else {
		// TODO: this dependency on the repositories package should be removed
		question.WithPrevious(repositories.DefaultProfileDir(), false)
	}
	config.SetProfilesDir(question.Ask(p))
}

func configOllama(p *output.Printer, config *config.Config) {
	q := output.NewQuestion("What is your Ollama host address?")
	if config.Ollama.Host != "" {
		q.WithPrevious(config.Ollama.Host, false)
	}
	config.Ollama.Host = q.Ask(p)

	q = output.NewQuestion("What Ollama model would you like to use?")
	if config.Ollama.Model != "" {
		q.WithPrevious(config.Ollama.Model, false)
	}
	config.Ollama.Model = q.Ask(p)
}
