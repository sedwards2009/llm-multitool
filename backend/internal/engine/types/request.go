package types

import "sedwards2009/llm-workbench/internal/data"

type Request struct {
	Prompt        string
	AppendFunc    func(string)
	CompleteFunc  func()
	SetStatusFunc func(status data.ResponseStatus)
	ModelSettings *data.ModelSettings
}
