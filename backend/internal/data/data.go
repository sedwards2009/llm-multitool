package data

import (
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/data/role"
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

type Session struct {
	ID                string         `json:"id"`
	CreationTimestamp string         `json:"creationTimestamp"`
	Title             string         `json:"title"`
	Prompt            string         `json:"prompt"`
	Responses         []*Response    `json:"responses"`
	ModelSettings     *ModelSettings `json:"modelSettings"`
}

type ModelSettings struct {
	ModelID    string `json:"modelId"`
	TemplateID string `json:"templateId"`
	PresetID   string `json:"presetId"`
}

type ModelOverview struct {
	Models []*Model `json:"models"`
}

type Model struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	EngineID        string
	InternalModelID string
}

type Response struct {
	ID                string                        `json:"id"`
	CreationTimestamp string                        `json:"creationTimestamp"`
	Status            responsestatus.ResponseStatus `json:"status"`
	Messages          []Message                     `json:"messages"`
}

type Message struct {
	ID   string    `json:"id"`
	Role role.Role `json:"role"`
	Text string    `json:"text"`
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
