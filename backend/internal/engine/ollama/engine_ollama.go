package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/engine/config"
	"sedwards2009/llm-workbench/internal/engine/types"
)

type OllamaEngineBackend struct {
	id        string
	config    *config.EngineBackendConfig
	isDefault bool
}

type modelList struct {
	Models []*model `json:"models"`
}

type model struct {
	Name string `json:"name"`
}

type generatePayload struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Options generateOptionsPayload `json:"options"`
}

type generateOptionsPayload struct {
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
}

type generateResponse struct {
	Response *string `json:"response"`
	Done     bool    `json:"done"`
}

func New(config *config.EngineBackendConfig) *OllamaEngineBackend {
	return &OllamaEngineBackend{
		id:        config.Name,
		config:    config,
		isDefault: config.Default,
	}
}

func (this *OllamaEngineBackend) ID() string {
	return this.id
}

func (this *OllamaEngineBackend) IsDefault() bool {
	return this.isDefault
}

func (this *OllamaEngineBackend) Process(work *types.Request, model *data.Model, preset *data.Preset) {
	log.Printf("OllamaEngineBackend process(): Starting request")
	work.SetStatusFunc(responsestatus.Running)
	defer work.CompleteFunc()

	payload := &generatePayload{
		Model:  model.InternalModelID,
		Prompt: work.Messages[0].Text,
		Options: generateOptionsPayload{
			Temperature: preset.Temperature,
			TopP:        preset.TopP,
		},
	}

	jsonData, _ := json.Marshal(payload)
	bodyBytes := bytes.NewBuffer(jsonData)
	url := *this.config.Address + "/api/generate"
	log.Printf("url: %s\n", url)
	resp, err := http.Post(url, "application/json", bodyBytes)
	if err != nil {
		log.Printf("OllamaEngineBackend Process(): CompletionStream error: %v\n", err)
		work.SetStatusFunc(responsestatus.Error)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("OllamaEngineBackend Process(): Error: HTTP request failed with status code: %d\n", resp.StatusCode)
		work.SetStatusFunc(responsestatus.Error)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		response := &generateResponse{}
		err := json.Unmarshal([]byte(line), &response)
		if err != nil {
			log.Printf("OllamaEngineBackend Process(): CompletionStream error: %v\n", err)
			work.SetStatusFunc(responsestatus.Error)
			return
		}
		if response.Done {
			break
		}
		if !work.AppendFunc(*response.Response) {
			break
		}
	}

	// Check for errors that may have occurred during scanning.
	if err := scanner.Err(); err != nil {
		log.Printf("OllamaEngineBackend Process(): CompletionStream error: %v\n", err)
		work.SetStatusFunc(responsestatus.Error)
		return
	}

	work.SetStatusFunc(responsestatus.Done)
	log.Printf("OllamaEngineBackend process(): CompletionStream completed")
}

func (this *OllamaEngineBackend) ScanModels() []*data.Model {
	url := *this.config.Address + "/api/tags"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("OllamaEngineBackend ScanModels(): Error: %v\n", err)
		return []*data.Model{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("OllamaEngineBackend ScanModels(): Error: HTTP request failed with status code: %d\n", resp.StatusCode)
		return []*data.Model{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("OllamaEngineBackend ScanModels(): Error: %v\n", err)
		return []*data.Model{}
	}

	modelList := &modelList{}
	if err := json.Unmarshal(body, &modelList); err != nil {
		log.Printf("OllamaEngineBackend ScanModels(): Error: %v\n", err)
		return []*data.Model{}
	}

	result := []*data.Model{}
	for _, modelInfo := range modelList.Models {
		result = append(result, &data.Model{
			ID:              this.id + "_" + modelInfo.Name,
			Name:            this.id + " - " + modelInfo.Name,
			EngineID:        this.id,
			InternalModelID: modelInfo.Name,
		})
	}
	return result
}
