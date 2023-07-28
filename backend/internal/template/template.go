package template

import (
	"sedwards2009/llm-workbench/internal/data"
	"strings"

	"github.com/bobg/go-generics/v2/slices"
)

type Templates struct {
	templates []*data.Template
}

const PROMPT_PARAM = "{{prompt}}"

func NewTemplates() *Templates {
	return &Templates{
		templates: []*data.Template{
			{
				ID:             "9e8df77f-c9c8-4683-995f-d744376901b5",
				Name:           "Instruct",
				TemplateString: PROMPT_PARAM,
			},

			{
				ID:             "4c080472-0aee-42cc-a100-e0fbb845f5a0",
				Name:           "Summerize text",
				TemplateString: "Summarize the following passage in a concise manner. Only give the summary.:\n\n" + PROMPT_PARAM,
			},
			{
				ID:             "91ccfaa5-5c46-4fee-bed4-37ca902789de",
				Name:           "Translate to Dutch",
				TemplateString: "Translate the following passage to Dutch. Only give the translation.:\n\n" + PROMPT_PARAM,
			},
			{
				ID:             "3b43850e-cc14-4516-85b9-4ba77e5939dd",
				Name:           "Translate to French",
				TemplateString: "Translate the following passage to French. Only give the translation.:\n\n" + PROMPT_PARAM,
			},
			{
				ID:             "e60c5bc8-132b-4d13-aab9-17043b234819",
				Name:           "Proofread and correct errors",
				TemplateString: "Proofread and edit the following text for any errors, typos, or mistakes:\n\n" + PROMPT_PARAM,
			},
		},
	}
}

func (this *Templates) TemplateOverview() *data.TemplateOverview {
	return &data.TemplateOverview{
		Templates: this.templates[:],
	}
}

func (this *Templates) ApplyTemplate(templateID string, promptText string) string {
	matches := slices.Filter(this.templates, func(t *data.Template) bool {
		return t.ID == templateID
	})
	if len(matches) == 0 {
		return promptText
	}
	return strings.Replace(matches[0].TemplateString, "{{prompt}}", promptText, -1)
}
