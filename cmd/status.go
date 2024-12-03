package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/core/config"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewStatusCmd(configRepo config.IConfigRepository, p *output.Printer) *cobra.Command {
	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Shows the current status of gennie",
		Long:  `Use it to check the current status of ginnie. You can check the model, profile and more!`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := configRepo.Load()
			if err != nil {
				ExitWithError(err)
			}

			p.PrintLine(output.Yellow)
			p.Print("API Keys", output.Cyan)
			apiKeyStatus("Open AI API Key", config.APIKeys.OpenAiApiKey, p)
			apiKeyStatus("Anthropic API Key", config.APIKeys.AnthropicApiKey, p)
			apiKeyStatus("Maritaca API Key", config.APIKeys.MaritacaApiKey, p)
			apiKeyStatus("Groq API Key", config.APIKeys.GroqApiKey, p)

			p.PrintLine(output.Yellow)
			p.Print("Ollama", output.Cyan)
			p.Print(fmt.Sprintf("Host: %s", config.Ollama.Host), output.Gray)
			p.Print(fmt.Sprintf("Model: %s", config.Ollama.Model), output.Gray)

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profiles path: %s", config.ProfilesDirPath), output.Gray)
			p.Print(fmt.Sprintf("Cache saved at: %s", config.ConversationCacheDir), output.Gray)
			p.PrintLine(output.Yellow)
		},
	}

	return cmdStatus
}

func apiKeyStatus(text string, apiKey string, p *output.Printer) {
	if apiKey != "" {
		p.Print(fmt.Sprintf("%s: Set", text), output.Green)
	} else {
		p.Print(fmt.Sprintf("%s: Not Set", text), output.Red)
	}
}
