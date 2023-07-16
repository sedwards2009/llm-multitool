package engine

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"sedwards2009/llm-workbench/internal/data"

	openai "github.com/sashabaranov/go-openai"
)

const ENGINE_NAME = "openai"

func processOpenAI(work *enqueueWorkPayload) {
	log.Printf("processOpenAI(): Starting request")
	work.setStatusFunc(data.ResponseStatus_Running)

	c := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 200,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: work.prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("processOpenAI(): ChatCompletionStream error: %v\n", err)
		work.completeFunc()
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
			work.setStatusFunc(data.ResponseStatus_Error)
			break
		}
		work.appendFunc(response.Choices[0].Delta.Content)
	}
	work.setStatusFunc(data.ResponseStatus_Done)
	log.Printf("processOpenAI(): ChatCompletionStream completed")
	work.completeFunc()
}

func scanModelsOpenAI() []*data.Model {
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
