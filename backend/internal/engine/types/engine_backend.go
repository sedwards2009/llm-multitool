package types

import "sedwards2009/llm-multitool/internal/data"

type EngineBackend interface {
	ID() string
	ScanModels() []*data.Model
	Process(work *Request, model *data.Model, preset *data.Preset)
}
