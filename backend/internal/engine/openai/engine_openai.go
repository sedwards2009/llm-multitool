package openai

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/engine/types"

	openai "github.com/sashabaranov/go-openai"
)

const ENGINE_NAME = "openai"

func NewEngineBackend() types.EngineBackend {
	return types.EngineBackend{
		ID:         ENGINE_NAME,
		ScanModels: scanModels,
		Process:    process,
	}
}

func process(work *types.Request, model *data.Model) {
	log.Printf("processOpenAI(): Starting request")
	work.SetStatusFunc(data.ResponseStatus_Running)

	c := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
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
			work.SetStatusFunc(data.ResponseStatus_Error)
			break
		}
		work.AppendFunc(response.Choices[0].Delta.Content)
	}
	work.SetStatusFunc(data.ResponseStatus_Done)
	log.Printf("processOpenAI(): ChatCompletionStream completed")
	work.CompleteFunc()
}

func scanModels() []*data.Model {
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
