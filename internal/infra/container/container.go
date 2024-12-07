package container

import (
	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/profile"
	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/infra/apiclient"
	"github.com/robertoseba/gennie/internal/infra/repositories"
)

type Container struct {
	conversationRepository conversation.IConversationRepository
	profileRepository      profile.IProfileRepository
	configRepository       config.IConfigRepository
	apiClient              models.IApiClient
	config                 *config.Config
}

func NewContainer() *Container {
	configDir, err := repositories.CreateConfigDir()
	if err != nil {
		panic(err)
	}
	configRepo := repositories.NewConfigRepository(configDir)
	config, err := configRepo.Load()
	if err != nil {
		panic(err)
	}

	container := &Container{
		config:                 config,
		configRepository:       configRepo,
		profileRepository:      repositories.NewProfileRepository(config.ProfilesDirPath),
		conversationRepository: repositories.NewConversationRepository(config.ConversationCacheDir),
		apiClient:              apiclient.NewApiClient(config.HttpTimeout),
	}

	return container
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetConfigRepository() config.IConfigRepository {
	return c.configRepository
}

// SERVICES
func (c *Container) GetCompleteService() *usecases.GetAnswerService {
	return usecases.NewGetAnswerService(
		c.conversationRepository,
		c.profileRepository,
		c.apiClient,
		c.config,
	)
}

func (c *Container) GetSelectModelService() *usecases.SelectModelService {
	return usecases.NewSelectModelService(c.conversationRepository)
}

func (c *Container) GetSelectProfileService() *usecases.SelectProfileService {
	return usecases.NewSelectProfileService(c.profileRepository, c.conversationRepository)
}

func (c *Container) GetExportConversationService() *usecases.ExportConversationService {
	return usecases.NewExportConversationService(c.conversationRepository)
}
