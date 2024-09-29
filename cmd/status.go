package cmd

import (
	"fmt"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/models"
	output "github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewStatusCmd(c *cache.Cache, p *output.Printer) *cobra.Command {

	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Shows the current status of gennie",
		Long:  `Use it to check the current status of ginnie. You can check the model, profile and more!`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if c.Model == "" || c.Profile == nil {
				ExitWithError(fmt.Errorf("Gennie hasn't been configured yet.\nPlease configure model and profile first, or just ask a question to get started with the default."))
			}

			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Model: %s ", models.ModelEnum(c.Model)), output.Cyan)
			p.PrintLine(output.Yellow)
			p.Print(fmt.Sprintf("Profile name: %s", c.Profile.Name), output.Gray)
			p.Print(fmt.Sprintf("Profile description: %s", c.Profile.Data), output.Gray)
			p.Print("", "")
			p.Print(fmt.Sprintf("Cache saved at: %s", c.FilePath), output.Gray)
			p.PrintLine(output.Yellow)
		},
	}

	return cmdStatus
}
