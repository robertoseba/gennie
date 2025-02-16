package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/spf13/cobra"
)

func NewConversationCmd(convService *usecases.ConversationService, p *output.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversation [command]",
		Short: "Manages conversations. Used to load/export conversations",
	}

	cmd.AddCommand(newExportConvCmd(convService, p))
	cmd.AddCommand(newLoadConversationCmd(convService, p))
	cmd.AddCommand(newLastConversationCmd(convService, p))
	return cmd
}

func newExportConvCmd(convService *usecases.ConversationService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "save [filename]",
		Short: "Saves the active conversation to a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Print(fmt.Sprintf("Saving conversation to: %s...", args[0]), output.Gray)
			err := convService.SaveTo(args[0])
			if err != nil {
				return err
			}
			p.Print("Conversation saved successfully", output.Green)
			return nil
		},
	}
}

func newLoadConversationCmd(convService *usecases.ConversationService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "load [filename]",
		Short: "Loads a conversation from a file and sets as active",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Print(fmt.Sprintf("Loading conversation from: %s...", args[0]), output.Gray)
			err := convService.LoadFrom(args[0])
			if err != nil {
				return err
			}
			p.Print("Conversation loaded successfully", output.Green)
			return nil
		},
	}
}

func newLastConversationCmd(convService *usecases.ConversationService, p *output.Printer) *cobra.Command {
	return &cobra.Command{
		Use:   "last",
		Short: "Shows the last conversation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			conv, err := convService.LastConversation()
			if err != nil {
				return err
			}
			if conv == nil {
				p.Print("No conversation could be found", output.Red)
				return nil
			}

			response, err := json.MarshalIndent(conv, "", "  ")
			if err != nil {
				return err
			}

			p.Print(string(response), output.Green)
			return nil
		},
	}
}
