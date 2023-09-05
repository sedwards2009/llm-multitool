package template

import (
	"fmt"
	"sedwards2009/llm-workbench/internal/data"
	"strings"

	"github.com/bobg/go-generics/v2/slices"
)

type Templates struct {
	templates []*data.Template
}

const PROMPT_PARAM = "{{prompt}}"

const TITLE_LENGTH = 40

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
	template := this.getTemplateByID(templateID)
	if template == nil {
		return promptText
	}

	return strings.Replace(template.TemplateString, "{{prompt}}", promptText, -1)
}

func (this *Templates) getTemplateByID(templateID string) *data.Template {
	matches := slices.Filter(this.templates, func(t *data.Template) bool {
		return t.ID == templateID
	})
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}

func (this *Templates) MakeTitle(templateID string, promptText string) string {
	lines := strings.Split(promptText, "\n")
	firstLine := lines[0]
	if len(firstLine) > TITLE_LENGTH {
		firstLine = lines[0][:TITLE_LENGTH]
	}

	template := this.getTemplateByID(templateID)
	if template == nil {
		return firstLine
	}

	return fmt.Sprintf("%s - %s", template.Name, firstLine)
}
