package oobabooga

import (
	"context"
	"errors"
	"io"
	"log"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/engine/types"

	"github.com/sashabaranov/go-openai"
)

const ENGINE_NAME = "oobabooga"

func NewEngineBackend() types.EngineBackend {
	return types.EngineBackend{
		ID:         ENGINE_NAME,
		ScanModels: scanModels,
		Process:    process,
	}
}

func process(work *types.Request, model *data.Model) {
	log.Printf("Process Oobabooga(): Starting request")
	work.SetStatusFunc(data.ResponseStatus_Running)

	c := openai.NewClientWithConfig(config())
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     model.InternalModelID,
		MaxTokens: 200,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: work.Prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("Process Oobabooga(): ChatCompletionStream error: %v\n", err)
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
			log.Printf("Process Oobabooga(): ChatCompletionStream error: %v\n", err)
			work.SetStatusFunc(data.ResponseStatus_Error)
			break
		}
		work.AppendFunc(response.Choices[0].Delta.Content)
	}
	work.SetStatusFunc(data.ResponseStatus_Done)
	log.Printf("Process Oobabooga(): ChatCompletionStream completed")
	work.CompleteFunc()
}

func config() openai.ClientConfig {
	config := openai.DefaultConfig("")
	config.BaseURL = "http://127.0.0.1:5001/v1"
	config.APIType = openai.APITypeOpenAI
	config.OrgID = ""
	return config
}

func scanModels() []*data.Model {
	c := openai.NewClientWithConfig(config())
	ctx := context.Background()

	result := []*data.Model{}

	modelList, err := c.ListModels(ctx)
	if err != nil {
		log.Printf("ScanModels Oobabooga(): Error: %v\n", err)
		return []*data.Model{}
	}

	for _, modelInfo := range modelList.Models {
		if modelInfo.Object != "model" {
			continue
		}
		result = append(result, &data.Model{
			ID:              "oobabooga_" + modelInfo.ID,
			Name:            "Oobabooga - " + modelInfo.ID,
			Engine:          ENGINE_NAME,
			InternalModelID: modelInfo.ID,
		})
		break
		// We only take the first one because Oobabooga doesn't
		// support loading different models on the fly.
	}
	return result
}
