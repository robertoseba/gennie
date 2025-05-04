package container

import (
	"net/http"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/profile"
	"github.com/robertoseba/gennie/internal/core/usecases"
	"github.com/robertoseba/gennie/internal/core/usecases/complete"
	"github.com/robertoseba/gennie/internal/infra/repositories"
)

type Container struct {
	conversationRepository conversation.IConversationRepository
	profileRepository      profile.IProfileRepository
	configRepository       config.IConfigRepository
	httpClient             *http.Client
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
		httpClient:             &http.Client{Timeout: config.HttpTimeout},
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
func (c *Container) GetCompleteService() *complete.CompleteService {
	return complete.NewCompleteService(
		c.conversationRepository,
		c.profileRepository,
		c.httpClient,
		c.config,
	)
}

func (c *Container) GetSelectModelService() *usecases.SelectModelService {
	return usecases.NewSelectModelService(c.conversationRepository)
}

func (c *Container) GetSelectProfileService() *usecases.SelectProfileService {
	return usecases.NewSelectProfileService(c.profileRepository, c.conversationRepository)
}

func (c *Container) GetExportConversationService() *usecases.ConversationService {
	return usecases.NewConversationService(c.conversationRepository)
}
