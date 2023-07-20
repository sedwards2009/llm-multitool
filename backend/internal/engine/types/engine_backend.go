package types

import "sedwards2009/llm-workbench/internal/data"

type EngineBackend struct {
	ID         string
	ScanModels func() []*data.Model
	Process    func(work *Request, model *data.Model)
}
