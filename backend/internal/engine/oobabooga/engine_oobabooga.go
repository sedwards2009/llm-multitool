package oobabooga

import (
	"context"
	"errors"
	"io"
	"log"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/engine/request"

	"github.com/sashabaranov/go-openai"
)

const ENGINE_NAME = "oobabooga"

func Process(work *request.Request) {
	log.Printf("Process Oobabooga(): Starting request")
	work.SetStatusFunc(data.ResponseStatus_Running)

	c := openai.NewClient("")
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
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

func ScanModels() []*data.Model {
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
			ID:              "oogabooga_" + modelInfo.ID,
			Name:            "Oogabooga - " + modelInfo.ID,
			Engine:          ENGINE_NAME,
			InternalModelID: modelInfo.ID,
		})
	}
	return result
}
