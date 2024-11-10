package client

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AuthEndpoint string `yaml:"authEndpoint,omitempty"`
	APIEndpoint  string `yaml:"apiEndpoint,omitempty"`
	AccessToken  string `yaml:"accessToken,omitempty"`
	RefreshToken string `yaml:"refreshToken,omitempty"`
	OrgName      string `yaml:"orgName,omitempty"`
}

func ReadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	var config Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &config, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}
