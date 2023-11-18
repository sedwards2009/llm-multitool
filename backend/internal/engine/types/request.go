package types

import (
	"sedwards2009/llm-multitool/internal/data"
	"sedwards2009/llm-multitool/internal/data/responsestatus"
)

type Request struct {
	Messages      []data.Message
	AppendFunc    func(string) bool
	CompleteFunc  func()
	SetStatusFunc func(status responsestatus.ResponseStatus)
	ModelSettings *data.ModelSettings
}
