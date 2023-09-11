package types

import "sedwards2009/llm-workbench/internal/data"

type EngineBackend interface {
	ID() string
	IsDefault() bool
	ScanModels() []*data.Model
	Process(work *Request, model *data.Model, preset *data.Preset)
}
