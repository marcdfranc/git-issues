package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git-issues/domain"
)

type Service struct {
	config *domain.Config
}

func New(config *domain.Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) MakeGitHubRequest(method, url string, data interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var reqBody []byte
	if data != nil {
		var err error
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("encoding error: %v", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error on request create: %v", err)
	}

	req.Header.Set("Authorization", "token "+s.config.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro on read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		var errorResponse struct {
			Message string `json:"message"`
		}
		json.Unmarshal(body, &errorResponse)
		return nil, fmt.Errorf("GitHub api error (%d): %s", resp.StatusCode, errorResponse.Message)
	}

	return body, nil
}
