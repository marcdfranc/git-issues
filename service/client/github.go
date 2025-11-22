package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"git-issues/domain"
)

var (
	errStr = "error on MakeRequest: %s"
)

type GitHubClient interface {
	MakeRequest(method, url string, data *domain.Issue) ([]byte, error)
}

type Service struct {
	config *domain.Config
}

func New(config *domain.Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) MakeRequest(method, url string, data *domain.Issue) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var reqBody []byte
	if data != nil {
		var err error
		reqBody, err = json.Marshal(data)
		if err != nil {
			err = fmt.Errorf(errStr, err)
			return nil, errors.Join(err, domain.ErrEncoding)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		err = fmt.Errorf(errStr, err)
		return nil, errors.Join(err, domain.ErrCreateRequest)
	}

	req.Header.Set("Authorization", "token "+s.config.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf(errStr, err)
		return nil, errors.Join(domain.ErrRequest, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(errStr, err)
	}

	if resp.StatusCode >= 400 {
		var errorResponse struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			err = fmt.Errorf(errStr, err)
			return nil, errors.Join(err, domain.ErrRequest)
		}
		err = fmt.Errorf("GitHub api error Status:%d\n response error: %s", resp.StatusCode, errorResponse.Message)
		return nil, errors.Join(err, domain.ErrApi)
	}

	return body, nil
}
