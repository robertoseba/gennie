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

			config.OpenAiApiKey = askKey(p, "OpenAI", config.OpenAiApiKey)
			config.AnthropicApiKey = askKey(p, "Anthropic", config.AnthropicApiKey)
			config.MaritacaApiKey = askKey(p, "Maritaca", config.MaritacaApiKey)
			config.GroqApiKey = askKey(p, "Groq", config.GroqApiKey)

			config.ProfilesPath = askProfileFolder(p, config.ProfilesPath)

			config.OllamaHost = p.Ask("Enter your Ollama host (ex: http://localhost:11434): ", output.Cyan)
			config.OllamaModel = p.Ask("Enter your Ollama model (ex: qwen2.5:0.5b): ", output.Cyan)
			storage.SetConfig(config)

			refreshProfiles(storage)

			return nil
		},
	}

	return clearCmd
}

func askKey(p *output.Printer, key string, previousValue string) string {
	if previousValue != "" {
		newValue := p.Ask(fmt.Sprintf("Enter your new %s API Key (or press ENTER keep the old one): ", key), output.Cyan)
		if newValue != "" {
			return newValue
		}
		return previousValue
	}
	return p.Ask(fmt.Sprintf("Enter your %s API Key: ", key), output.Cyan)
}

func askProfileFolder(p *output.Printer, previousValue string) string {
	if previousValue != "" {
		newValue := p.Ask(fmt.Sprintf("Enter your new profiles folder path (or press ENTER keep %s): ", previousValue), output.Cyan)
		if newValue != "" {
			return newValue
		}
		return previousValue
	}

	defaultPath := profile.DefaultProfilePath()
	profileFolder := p.Ask(fmt.Sprintf("Enter the path to your profiles folder or press ENTER to use Default(%s): ", defaultPath), output.Cyan)

	if profileFolder == "" {
		profileFolder = defaultPath
	}

	return profileFolder
}
