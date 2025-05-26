package issue

import (
	"encoding/json"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
)

type ListIssue interface {
	List() ([]domain.Issue, error)
}

type ListFeature struct {
	config *domain.Config
	client client.GitHubClient
}

func NewList(config *domain.Config, client client.GitHubClient) *ListFeature {
	return &ListFeature{
		config: config,
		client: client,
	}
}

func (f *ListFeature) List() ([]domain.Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", f.config.APIBaseURL, f.config.Owner, f.config.Repo)

	response, err := f.client.MakeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	issues := []domain.Issue{}
	if err = json.Unmarshal(response, &issues); err != nil {
		return nil, errProcessing
	}

	return issues, nil
}
