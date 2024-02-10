package data

import (
	"sedwards2009/llm-multitool/internal/data/responsestatus"
	"sedwards2009/llm-multitool/internal/data/role"
)

type SessionOverview struct {
	SessionSummaries []*SessionSummary `json:"sessionSummaries"`
}

type SessionSummary struct {
	ID                string `json:"id"`
	CreationTimestamp string `json:"creationTimestamp"`
	Title             string `json:"title"`
}

type Root struct {
	Sessions []Session `json:"sessions"`
}

type AttachedFile struct {
	Filename         string `json:"filename"`
	MimeType         string `json:"mimeType"`
	OriginalFilename string `json:"originalFilename"`
}

type Session struct {
	ID                string          `json:"id"`
	CreationTimestamp string          `json:"creationTimestamp"`
	Title             string          `json:"title"`
	Prompt            string          `json:"prompt"`
	AttachedFiles     []*AttachedFile `json:"attachedFiles"`
	Responses         []*Response     `json:"responses"`
	ModelSettings     *ModelSettings  `json:"modelSettings"`
}

type ModelSettings struct {
	ModelID    string `json:"modelId"`
	TemplateID string `json:"templateId"`
	PresetID   string `json:"presetId"`
}

type ModelSettingsSnapshot struct {
	ModelSettings
	ModelName    string `json:"modelName"`
	TemplateName string `json:"templateName"`
	PresetName   string `json:"presetName"`
}

type ModelOverview struct {
	Models []*Model `json:"models"`
}

type Model struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	EngineID         string
	InternalModelID  string
	SupportsContinue bool `json:"supportsContinue"`
	SupportsReply    bool `json:"supportsReply"`
}

type Response struct {
	ID                    string                        `json:"id"`
	CreationTimestamp     string                        `json:"creationTimestamp"`
	Status                responsestatus.ResponseStatus `json:"status"`
	Messages              []Message                     `json:"messages"`
	ModelSettingsSnapshot *ModelSettingsSnapshot        `json:"modelSettingsSnapshot"`
}

type Message struct {
	ID            string          `json:"id"`
	Role          role.Role       `json:"role"`
	Text          string          `json:"text"`
	AttachedFiles []*AttachedFile `json:"attachedFiles"`
}

type Template struct {
	ID             string `json:"id" yaml:"id"`
	Name           string `json:"name" yaml:"name"`
	TemplateString string `json:"templateString" yaml:"template_string"`
	Default        bool   `yaml:"default,omitempty"`
}

type TemplateOverview struct {
	Templates []*Template `json:"templates"`
}

type Preset struct {
	ID          string  `json:"id" yaml:"id"`
	Name        string  `json:"name" yaml:"name"`
	Temperature float32 `yaml:"temperature"`
	TopP        float32 `yaml:"top_p"`
	Default     bool    `yaml:"default,omitempty"`
}

type PresetOverview struct {
	Presets []*Preset `json:"presets"`
}
