package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Auth struct {
	PAT    string `json:"pat"`
	Secret string `json:"secret"`
}

type Rule struct {
	Consume    string   `json:"consume"`
	Produce    string   `json:"produce"`
	Formatters []string `json:"formatters"`
}

type Config struct {
	Port  int    `json:"port"`
	Auth  Auth   `json:"auth"`
	Rules []Rule `json:"rules"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}
	result := &Config{}
	if err = yaml.Unmarshal(data, result); err != nil {
		return nil, err
	}
	if result.Port == 0 {
		result.Port = 20316
	}
	return result, nil
}
