package types

import (
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
)

type Request struct {
	Messages      []data.Message
	AppendFunc    func(string)
	CompleteFunc  func()
	SetStatusFunc func(status responsestatus.ResponseStatus)
	ModelSettings *data.ModelSettings
}
