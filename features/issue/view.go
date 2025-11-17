package issue

import (
	"encoding/json"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
)

type ViewIssue interface {
	View(issueNumber int) (*domain.Issue, error)
}

type ViewFeature struct {
	config *domain.Config
	client client.GitHubClient
}

func NewView(config *domain.Config, client client.GitHubClient) *ViewFeature {
	return &ViewFeature{
		config: config,
		client: client,
	}
}

func (f *ViewFeature) View(issueNumber int) (*domain.Issue, error) {
	if issueNumber == 0 {
		return nil, errNumberIsRequered
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", f.config.APIBaseURL, f.config.Owner, f.config.Repo, issueNumber)

	response, err := f.client.MakeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	issue := &domain.Issue{}
	if err = json.Unmarshal(response, &issue); err != nil {
		return nil, errProcessing
	}
	return issue, nil
}
