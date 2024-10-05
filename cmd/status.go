package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/common"
	"github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewStatusCmd(persistence common.IPersistence, p *output.Printer) *cobra.Command {

	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Shows the current status of gennie",
		Long:  `Use it to check the current status of ginnie. You can check the model, profile and more!`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := persistence.GetConfig()

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Model: %s ", models.ModelEnum(config.CurrModelSlug)), output.Cyan)
			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profile name: %s", config.CurrProfile.Name), output.Gray)
			p.Print(fmt.Sprintf("Profile description: %s", config.CurrProfile.Data), output.Gray)
			p.Print("", "")
			p.Print(fmt.Sprintf("Cache saved at: %s", persistence.GetCacheFilePath()), output.Gray)
			p.PrintLine(output.Yellow)
		},
	}

	return cmdStatus
}
