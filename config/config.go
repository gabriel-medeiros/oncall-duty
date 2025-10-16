package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MaxDiff      int `json:"maxDiff"`
	DescansoDias int `json:"descansoDias"`
}

func Load(filename string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
