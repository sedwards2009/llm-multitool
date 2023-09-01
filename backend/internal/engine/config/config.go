package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type EngineBackendConfig struct {
	Name         string  `yaml:"name"`
	Address      *string `yaml:"address"`
	ApiTokenFrom *string `yaml:"api_token_from"`
	ApiToken     string
	Models       *[]string `yaml:"models"`
}

func ReadConfigFile(file string) ([]*EngineBackendConfig, error) {
	backendConfigs := &[]*EngineBackendConfig{}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	if err := yaml.Unmarshal(f, backendConfigs); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}

	loadApiTokens(backendConfigs, path.Dir(file))
	return *backendConfigs, nil
}

func loadApiTokens(backendConfigs *[]*EngineBackendConfig, basePath string) {
	for _, config := range *backendConfigs {
		if config.ApiTokenFrom != nil {
			apiTokenPath := path.Join(basePath, *config.ApiTokenFrom)
			content, err := os.ReadFile(apiTokenPath)
			if err != nil {
				fmt.Printf("Error opening api_token_from file '%s':\n%s", apiTokenPath, err.Error())
			} else {
				config.ApiToken = strings.TrimSpace(string(content))
			}
		} else {
			config.ApiToken = ""
		}
	}
}
