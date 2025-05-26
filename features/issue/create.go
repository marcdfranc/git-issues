package issue

import (
	"encoding/json"
	"errors"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
)

type Create interface {
	Create() error
}

type CreateFeature struct {
	client client.GitHubClient
	editor editor.Editor
	config *domain.Config
}

func NewCreate(config *domain.Config, editor editor.Editor, client client.GitHubClient) *CreateFeature {
	return &CreateFeature{
		config: config,
		editor: editor,
		client: client,
	}
}

func (f *CreateFeature) Create() (string, error) {
	issue := &domain.Issue{}

	err := f.editor.GetIssueContentFromEditor(issue)
	if err != nil {
		return "", errors.Join(err, domain.ErrEditor)
	}

	if issue.Title == "" {
		return "", errTitleRequired
	}

	if issue.Body == "" {
		return "", errBodyRequired
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues", f.config.APIBaseURL, f.config.Owner, f.config.Repo)
	response, err := f.client.MakeRequest("POST", url, issue)
	if err != nil {
		return "", errors.Join(err, errCreate)
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return "", errProcessing
	}

	return fmt.Sprintf("Issue created with success!\nNumber: %v\nURL: %v\n", result["number"], result["html_url"]), nil
}
