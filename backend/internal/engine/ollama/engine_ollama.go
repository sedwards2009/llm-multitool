package ollama

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sedwards2009/llm-multitool/internal/data"
	"sedwards2009/llm-multitool/internal/data/responsestatus"
	"sedwards2009/llm-multitool/internal/data/role"
	"sedwards2009/llm-multitool/internal/engine/config"
	"sedwards2009/llm-multitool/internal/engine/types"

	"github.com/bobg/go-generics/v2/slices"
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

type chatMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type chatPayload struct {
	Model    string         `json:"model"`
	Messages []chatMessage  `json:"messages"`
	Options  optionsPayload `json:"options"`
}

type chatResponse struct {
	Message *chatMessage `json:"message,omitempty"`

	Done bool `json:"done"`
}

type optionsPayload struct {
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
}

func New(config *config.EngineBackendConfig) *OllamaEngineBackend {
	return &OllamaEngineBackend{
		id:     config.Name,
		config: config,
	}
}

func (this *OllamaEngineBackend) ID() string {
	return this.id
}

func (this *OllamaEngineBackend) Process(work *types.Request, model *data.Model, preset *data.Preset) {
	log.Printf("OllamaEngineBackend process(): Starting request")
	work.SetStatusFunc(responsestatus.Running)
	defer work.CompleteFunc()
	log.Printf("OllamaEngineBackend Process(): Temperature: %f, TopP: %f\n", preset.Temperature, preset.TopP)
	previousMessages := work.Messages[0 : len(work.Messages)-1]
	payload := &chatPayload{
		Model: model.InternalModelID,
		Messages: slices.Map(previousMessages, func(m data.Message) chatMessage {
			mRole := "user"
			if m.Role == role.Assistant {
				mRole = "assistant"
			}

			images := []string{}
			if len(m.AttachedFiles) != 0 {
				images = slices.Map(m.AttachedFiles, func(af *data.AttachedFile) string {
					return readFileBase64(filepath.Join(work.AttachedFilesPath, af.Filename))
				})
			}

			return chatMessage{
				Role:    mRole,
				Content: m.Text,
				Images:  images,
			}
		}),
		Options: optionsPayload{
			Temperature: preset.Temperature,
			TopP:        preset.TopP,
		},
	}

	jsonData, _ := json.Marshal(payload)
	bodyBytes := bytes.NewBuffer(jsonData)
	url := *this.config.Address + "/api/chat"
	resp, err := http.Post(url, "application/json", bodyBytes)
	if err != nil {
		log.Printf("OllamaEngineBackend Process(): ChatStream error: %v\n", err)
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
		response := &chatResponse{}
		err := json.Unmarshal([]byte(line), &response)
		if err != nil {
			log.Printf("OllamaEngineBackend Process(): ChatStream error: %v\n", err)
			work.SetStatusFunc(responsestatus.Error)
			return
		}
		if response.Done {
			break
		}
		if !work.AppendFunc(*&response.Message.Content) {
			break
		}
	}

	// Check for errors that may have occurred during scanning.
	if err := scanner.Err(); err != nil {
		log.Printf("OllamaEngineBackend Process(): ChatStream error: %v\n", err)
		work.SetStatusFunc(responsestatus.Error)
		return
	}

	work.SetStatusFunc(responsestatus.Done)
	log.Printf("OllamaEngineBackend process(): ChatStream completed")
}

func readFileBase64(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("OllamaEngineBackend: Error reading file %s. %v\n", filePath, err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(content)
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
			ID:               this.id + "_" + modelInfo.Name,
			Name:             this.id + " - " + modelInfo.Name,
			EngineID:         this.id,
			InternalModelID:  modelInfo.Name,
			SupportsContinue: false,
			SupportsReply:    true,
		})
	}
	return result
}
