package openai

import (
	"context"
	"errors"
	"io"
	"log"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/data/role"
	"sedwards2009/llm-workbench/internal/engine/config"
	"sedwards2009/llm-workbench/internal/engine/types"

	"github.com/bobg/go-generics/v2/slices"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAiEngineBackend struct {
	id        string
	config    *config.EngineBackendConfig
	isDefault bool
}

func New(config *config.EngineBackendConfig) *OpenAiEngineBackend {
	return &OpenAiEngineBackend{
		id:        config.Name,
		config:    config,
		isDefault: config.Default,
	}
}

func (this *OpenAiEngineBackend) formatApiConfig() openai.ClientConfig {
	apiConfig := openai.DefaultConfig(this.config.ApiToken)
	if this.config.Address != nil {
		apiConfig.BaseURL = *this.config.Address
	}
	apiConfig.APIType = openai.APITypeOpenAI
	apiConfig.OrgID = ""
	return apiConfig
}

func (this *OpenAiEngineBackend) ID() string {
	return this.id
}

func (this *OpenAiEngineBackend) IsDefault() bool {
	return this.isDefault
}

func (this *OpenAiEngineBackend) Process(work *types.Request, model *data.Model, preset *data.Preset) {
	log.Printf("OpenAiEngineBackend process(): Starting request")
	work.SetStatusFunc(responsestatus.Running)
	defer work.CompleteFunc()

	c := openai.NewClientWithConfig(this.formatApiConfig())
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     model.InternalModelID,
		MaxTokens: 500,
		Messages: slices.Map(work.Messages, func(m data.Message) openai.ChatCompletionMessage {
			openaiRole := openai.ChatMessageRoleUser
			if m.Role == role.Assistant {
				openaiRole = openai.ChatMessageRoleAssistant
			}
			return openai.ChatCompletionMessage{
				Role:    openaiRole,
				Content: m.Text,
			}
		}),
		Stream: true,

		Temperature: preset.Temperature,
		TopP:        preset.TopP,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("OpenAiEngineBackend process(): ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("OpenAiEngineBackend process(): ChatCompletionStream error: %v\n", err)
			work.SetStatusFunc(responsestatus.Error)
			break
		}
		if !work.AppendFunc(response.Choices[0].Delta.Content) {
			break
		}
	}
	work.SetStatusFunc(responsestatus.Done)
	log.Printf("OpenAiEngineBackend process(): ChatCompletionStream completed")
}

func (this *OpenAiEngineBackend) ScanModels() []*data.Model {
	c := openai.NewClientWithConfig(this.formatApiConfig())
	ctx := context.Background()

	result := []*data.Model{}

	modelList, err := c.ListModels(ctx)
	if err != nil {
		log.Printf("OpenAiEngineBackend ScanModels(): Error: %v\n", err)
		return []*data.Model{}
	}

	models := make(map[string]bool)
	isFilterModels := this.config.Models != nil
	if isFilterModels {
		for _, model := range *this.config.Models {
			models[model] = true
		}
	}

	for _, modelInfo := range modelList.Models {
		if modelInfo.Object != "model" {
			continue
		}

		_, hasModels := models[modelInfo.ID]
		if isFilterModels && !hasModels {
			continue
		}

		result = append(result, &data.Model{
			ID:               this.id + "_" + modelInfo.ID,
			Name:             this.id + " - " + modelInfo.ID,
			EngineID:         this.id,
			InternalModelID:  modelInfo.ID,
			SupportsContinue: true,
			SupportsReply:    true,
		})

		if this.config.Variant != nil && *this.config.Variant == config.VARIANT_OOBABOOGA {
			// We only take the first one because Oobabooga doesn't
			// support loading different models on the fly.
			break
		}
	}
	return result
}
