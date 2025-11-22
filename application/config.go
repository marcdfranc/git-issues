package application

import (
	"encoding/json"
	"os"

	"git-issues/domain"
)

func LoadConfig(filePath string) (*domain.Config, error) {
	data, err := os.ReadFile(filePath)
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
