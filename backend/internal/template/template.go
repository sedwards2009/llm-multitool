package template

import (
	"fmt"
	"os"
	"sedwards2009/llm-multitool/internal/data"
	"strings"

	"github.com/bobg/go-generics/v2/slices"
	"gopkg.in/yaml.v3"
)

type TemplateDatabase struct {
	templates []*data.Template
}

const PROMPT_PARAM = "{{prompt}}"

const TITLE_LENGTH = 40

func MakeTemplateDatabase(fileName string) (*TemplateDatabase, error) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Cannot read templates file '%s': %w", fileName, err)
	}
	return MakeTemplateDatabaseFromBytes(fileContents, fileName)
}

func MakeTemplateDatabaseFromBytes(yamlBytes []byte, fileName string) (*TemplateDatabase, error) {
	this := &TemplateDatabase{
		templates: make([]*data.Template, 0),
	}
	err := this.readTemplatesYamlString(yamlBytes, fileName)
	return this, err
}

func (this *TemplateDatabase) readTemplatesYamlString(yamlContent []byte, fileName string) error {
	if err := yaml.Unmarshal(yamlContent, &this.templates); err != nil {
		return fmt.Errorf("Cannot unmarshal config file '%s': %w", fileName, err)
	}
	return nil
}

func (this *TemplateDatabase) TemplateOverview() *data.TemplateOverview {
	return &data.TemplateOverview{
		Templates: this.templates[:],
	}
}

func (this *TemplateDatabase) DefaultID() string {
	for _, template := range this.templates {
		if template.Default {
			return template.ID
		}
	}
	return ""
}

func (this *TemplateDatabase) ApplyTemplate(templateID string, promptText string) string {
	template := this.getTemplateByID(templateID)
	if template == nil {
		return promptText
	}

	return strings.Replace(template.TemplateString, "{{prompt}}", promptText, -1)
}

func (this *TemplateDatabase) Get(templateID string) *data.Template {
	return this.getTemplateByID(templateID)
}

func (this *TemplateDatabase) getTemplateByID(templateID string) *data.Template {
	matches := slices.Filter(this.templates, func(t *data.Template) bool {
		return t.ID == templateID
	})
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}

func (this *TemplateDatabase) MakeTitle(templateID string, promptText string) string {
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
