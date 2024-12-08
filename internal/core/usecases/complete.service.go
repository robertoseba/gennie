package usecases

import (
	"os"
	"strings"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	apiClient        models.IApiClient
	config           *config.Config
}

type InputDTO struct {
	Question     string
	ProfileSlug  string
	Model        string
	IsFollowUp   bool
	AppendFile   string
	IsStreamable bool
}

func NewCompleteService(
	cr conversation.IConversationRepository,
	pr profile.IProfileRepository,
	apiClient models.IApiClient,
	config *config.Config) *CompleteService {
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		apiClient:        apiClient,
		config:           config,
	}
}

func (s *CompleteService) Execute(input *InputDTO) (<-chan models.StreamResponse, error) {
	var conv *conversation.Conversation
	var err error

	conv, err = s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}

	if input.ProfileSlug == "" {
		input.ProfileSlug = conv.ProfileSlug
	}

	if input.Model == "" {
		input.Model = conv.ModelSlug
	}

	if !input.IsFollowUp {
		conv = conversation.NewConversation(input.ProfileSlug, input.Model)
	}

	if input.AppendFile != "" {
		content, err := readFile(input.AppendFile)
		if err != nil {
			return nil, err
		}

		input.Question = input.Question + "\n" + content
	}

	if err := conv.NewQuestion(input.Question); err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.FindBySlug(input.ProfileSlug)
	if err != nil {
		return nil, err
	}

	if input.ProfileSlug != conv.ProfileSlug {
		conv.SetProfileTo(input.ProfileSlug)
	}

	model, err := models.NewModel(input.Model, s.apiClient, *s.config)
	if err != nil {
		return nil, err
	}

	if conv.ModelSlug != input.Model {
		conv.SetModelTo(input.Model)
	}

	outputChan := make(chan models.StreamResponse)
	go func() {
		defer close(outputChan)

		if !model.CanStream() {
			err = model.Complete(conv, profile.Data)
			if err != nil {
				outputChan <- models.StreamResponse{Err: err}
				return
			}
			outputChan <- models.StreamResponse{Data: conv.LastAnswer(), Err: nil}
			return
		}

		respChan, err := model.CompleteStreamable(conv, profile.Data)
		if err != nil {
			outputChan <- models.StreamResponse{Err: err}
			return
		}

		buf := strings.Builder{}
		for d := range respChan {
			outputChan <- d
			buf.WriteString(d.Data)
		}

		conv.AnswerLastQuestion(buf.String())
		err = s.conversationRepo.SaveAsActive(conv)
		if err != nil {
			outputChan <- models.StreamResponse{Err: err}
		}
	}()

	return outputChan, nil
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}
