package presets

import (
	"fmt"
	"log"
	"os"
	"sedwards2009/llm-workbench/internal/data"

	"gopkg.in/yaml.v3"
)

type PresetDatabase struct {
	presets []*data.Preset
}

func MakePresetDatabase(fileName string) (*PresetDatabase, error) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Cannot read presets file '%s': %w", fileName, err)
	}
	return MakePresentDatabaseFromBytes(fileContents, fileName)
}

func MakePresentDatabaseFromBytes(yamlBytes []byte, fileName string) (*PresetDatabase, error) {
	this := &PresetDatabase{
		presets: make([]*data.Preset, 0),
	}
	err := this.readPresentsYamlString(yamlBytes, fileName)
	return this, err
}

func (this *PresetDatabase) readPresentsYamlString(yamlContent []byte, fileName string) error {
	if err := yaml.Unmarshal(yamlContent, &this.presets); err != nil {
		return fmt.Errorf("Cannot unmarshal presets file '%s': %w", fileName, err)
	}
	return nil
}

func (this *PresetDatabase) PresetOverview() *data.PresetOverview {
	return &data.PresetOverview{
		Presets: this.presets,
	}
}

func (this *PresetDatabase) Exists(presetID string) bool {
	return this.Get(presetID) != nil
}

func (this *PresetDatabase) Get(presetID string) *data.Preset {
	for _, preset := range this.presets {
		if preset.ID == presetID {
			return preset
		}
	}
	log.Printf("PresetDatabase could not find preset %s", presetID)
	return nil
}

func (this *PresetDatabase) DefaultID() string {
	for _, preset := range this.presets {
		if preset.Default {
			return preset.ID
		}
	}
	return ""
}
