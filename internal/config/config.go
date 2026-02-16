package config

import (
	"encoding/json"
	"os"
)

type AppConfig map[string]string

func LoadConfig(path string) (AppConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg AppConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
