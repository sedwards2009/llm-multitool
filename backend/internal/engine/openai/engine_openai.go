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

const ENGINE_NAME = "openai"

type OpenAiEngineBackend struct {
	id     string
	config *config.EngineBackendConfig
}

func NewEngineBackend(config *config.EngineBackendConfig) OpenAiEngineBackend {
	return OpenAiEngineBackend{
		id:     ENGINE_NAME,
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
	return []*data.Model{
		{
			ID:              "openai.com_chatgpt3.5turbo",
			Name:            "OpenAI - ChatGPT 3.5 Turbo",
			Engine:          ENGINE_NAME,
			InternalModelID: openai.GPT3Dot5Turbo,
		},
		{
			ID:              "openai.com_gpt4",
			Name:            "OpenAI - GPT 4",
			Engine:          ENGINE_NAME,
			InternalModelID: openai.GPT4,
		},
	}
}
