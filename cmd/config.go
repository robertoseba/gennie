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

			config.OpenAiApiKey = p.Ask("Enter your Open AI API Key: ", output.Cyan)
			config.AnthropicApiKey = p.Ask("Enter your Anthopic API Key: ", output.Cyan)
			config.MaritacaApiKey = p.Ask("Enter your Maritaca API Key: ", output.Cyan)

			defaultPath := profile.DefaultProfilePath()
			profileFolder := p.Ask(fmt.Sprintf("Enter the path to your profiles folder or press ENTER to use Default(%s): ", defaultPath), output.Cyan)
			if profileFolder == "" {
				profileFolder = defaultPath
			}
			config.ProfilesPath = profileFolder

			storage.SetConfig(config)

			refreshProfiles(storage)

			return nil
		},
	}

	return clearCmd
}
