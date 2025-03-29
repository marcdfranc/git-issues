package application

import (
	"encoding/json"
	"os"

	"git-issues/domain"
)

func LoadConfig() (*domain.Config, error) {
	data, err := os.ReadFile(domain.ConfigFile)
	if err != nil {
		return nil, err
	}

	var config domain.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
