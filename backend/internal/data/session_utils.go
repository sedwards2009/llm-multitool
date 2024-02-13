package data

import "github.com/bobg/go-generics/v2/slices"

func (this *Session) GetAttachedFiles() []*AttachedFile {
	result := []*AttachedFile{}
	result = append(result, this.AttachedFiles...)

	for _, response := range this.Responses {
		for _, message := range response.Messages {
			result = append(result, message.AttachedFiles...)
		}
	}
	return result
}

func (this *Session) GetAttachedFileNames() []string {
	return slices.Map(this.GetAttachedFiles(), func(af *AttachedFile) string {
		return af.Filename
	})
}
