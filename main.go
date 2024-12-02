package main

import (
	_ "embed"

	"github.com/robertoseba/gennie/cmd"
	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/infra/apiclient"
	"github.com/robertoseba/gennie/internal/infra/repositories"
	"github.com/robertoseba/gennie/internal/output"
)

//go:embed version.txt
var version string

func main() {
	configDir, err := repositories.CreateConfigDir()
	if err != nil {
		panic(err)
	}
	configRepo := repositories.NewConfigRepository(configDir)
	config, err := configRepo.Load()
	if err != nil {
		panic(err)
	}

	//Repos
	conversationRepo := repositories.NewConversationRepository(config.ConversationCacheDir)
	profileRepo := repositories.NewProfileRepository(config.ProfilesDirPath)

	//ApiClient
	apiClient := apiclient.NewApiClient(config.HttpTimeout)

	//Cmds
	askCmd := usecases.NewGetAnswerService(conversationRepo, profileRepo, apiClient, *config)
	selectModelCmd := usecases.NewSelectModelService(conversationRepo)
	selectProfileCmd := usecases.NewSelectProfileService(profileRepo, conversationRepo)
	exportConversationCmd := usecases.NewExportConversationService(conversationRepo)

	//Output
	printer := output.NewPrinter(nil, nil)

	command := cmd.NewRootCmd(printer, version, config, configRepo, askCmd, selectModelCmd, selectProfileCmd, exportConversationCmd)

	if config.IsNew() {
		command.SetArgs([]string{"config"})
	}

	command.Execute()
}
