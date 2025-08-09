package issue

import (
	"encoding/json"
	"errors"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
)

type CloseIssue interface {
	Close(number int) error
}

type CloseFeature struct {
	config *domain.Config
	client client.GitHubClient
}

func NewClose(config *domain.Config, client client.GitHubClient) *CloseFeature {
	return &CloseFeature{
		config: config,
		client: client,
	}
}

func (f *CloseFeature) Close(number int) error {
	if number == 0 {
		return errNumberIsRequered
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", f.config.APIBaseURL, f.config.Owner, f.config.Repo, number)

	response, err := f.client.MakeRequest("GET", url, nil)
	if err != nil {
		return errNotFound
	}

	issue := &domain.Issue{}
	err = json.Unmarshal(response, issue)
	if err != nil {
		return errors.Join(errProcessing)
	}

	issue.State = "closed"

	response, err = f.client.MakeRequest("PATCH", url, issue)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(response, &issue); err != nil {
		return errProcessing
	}

	if issue.State != "closed" {
		return errClose
	}

	return nil
}
