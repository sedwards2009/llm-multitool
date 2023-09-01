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
	id     string
	config *config.EngineBackendConfig
}

func NewEngineBackend(config *config.EngineBackendConfig) OpenAiEngineBackend {
	return OpenAiEngineBackend{
		id:     config.Name,
		config: config,
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

func (this OpenAiEngineBackend) ID() string {
	return this.id
}

func (this OpenAiEngineBackend) Process(work *types.Request, model *data.Model) {
	log.Printf("processOpenAI(): Starting request")
	work.SetStatusFunc(responsestatus.Running)

	c := openai.NewClientWithConfig(this.formatApiConfig())
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     model.InternalModelID,
		MaxTokens: 200,
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
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("processOpenAI(): ChatCompletionStream error: %v\n", err)
		work.CompleteFunc()
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("processOpenAI(): ChatCompletionStream error: %v\n", err)
			work.SetStatusFunc(responsestatus.Error)
			break
		}
		work.AppendFunc(response.Choices[0].Delta.Content)
	}
	work.SetStatusFunc(responsestatus.Done)
	log.Printf("processOpenAI(): ChatCompletionStream completed")
	work.CompleteFunc()
}

func (this OpenAiEngineBackend) ScanModels() []*data.Model {
	c := openai.NewClientWithConfig(this.formatApiConfig())
	ctx := context.Background()

	result := []*data.Model{}

	modelList, err := c.ListModels(ctx)
	if err != nil {
		log.Printf("ScanModels: Error: %v\n", err)
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
			ID:              this.id + "_" + modelInfo.ID,
			Name:            this.id + " - " + modelInfo.ID,
			Engine:          this.id,
			InternalModelID: modelInfo.ID,
		})
		// break
		// We only take the first one because Oobabooga doesn't
		// support loading different models on the fly.
	}
	return result
}
