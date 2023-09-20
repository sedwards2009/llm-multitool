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

func MakePresetDatabase(file string) *PresetDatabase {
	this := &PresetDatabase{
		presets: make([]*data.Preset, 0),
	}
	err := this.readPresetsFile(file)
	if err != nil {
		log.Print(err)
	}
	return this
}

func (this *PresetDatabase) readPresetsFile(file string) error {
	f, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Cannot read presets file '%s': %w", file, err)
	}
	if err := yaml.Unmarshal(f, &this.presets); err != nil {
		return fmt.Errorf("Cannot unmarshal presets file '%s': %w", file, err)
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
