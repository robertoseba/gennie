package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewStatusCmd(storage common.IStorage, p *output.Printer) *cobra.Command {

	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Shows the current status of gennie",
		Long:  `Use it to check the current status of ginnie. You can check the model, profile and more!`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			profile := storage.GetCurrProfile()

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Model: %s ", models.ModelEnum(storage.GetCurrModelSlug())), output.Cyan)
			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profile name: %s", profile.Name), output.Gray)
			p.Print(fmt.Sprintf("Profile description: %s", profile.Data), output.Gray)

			config := storage.GetConfig()

			p.PrintLine(output.Yellow)
			p.Print("API Keys", output.Cyan)
			apiKeyStatus("Open AI API Key", config.OpenAiApiKey, p)
			apiKeyStatus("Anthropic API Key", config.AnthropicApiKey, p)
			apiKeyStatus("Maritaca API Key", config.MaritacaApiKey, p)

			p.PrintLine(output.Yellow)
			p.Print("Ollama", output.Cyan)
			p.Print(fmt.Sprintf("Host: %s", config.OllamaHost), output.Gray)
			p.Print(fmt.Sprintf("Model: %s", config.OllamaModel), output.Gray)

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profiles path: %s", config.ProfilesPath), output.Gray)
			p.Print(fmt.Sprintf("Cache saved at: %s", storage.GetStorageFilepath()), output.Gray)
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
