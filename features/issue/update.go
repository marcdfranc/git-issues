package issue

import (
	"encoding/json"
	"errors"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
)

type UpdateIssue interface {
	Update(number int) error
}

type UpdateFeature struct {
	config *domain.Config
	client client.GitHubClient
	editor editor.Editor
}

func NewUpdate(config *domain.Config, editor editor.Editor, client client.GitHubClient) *UpdateFeature {
	return &UpdateFeature{
		config: config,
		editor: editor,
		client: client,
	}
}

func (f *UpdateFeature) Update(number int) error {
	if number == 0 {
		return errNumberIsRequered
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", f.config.APIBaseURL, f.config.Owner, f.config.Repo, number)
	response, err := f.client.MakeRequest("GET", url, nil)
	if err != nil {
		return errNotFound
	}

	existingIssue := &domain.Issue{}
	err = json.Unmarshal(response, existingIssue)
	if err != nil {
		return errors.Join(errProcessing)
	}

	err = f.editor.GetIssueContentFromEditor(existingIssue)
	if err != nil {
		return errUpdate
	}

	response, err = f.client.MakeRequest("PATCH", url, existingIssue)
	if err != nil {
		return errors.Join(errUpdate, err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return errors.Join(errProcessing, err)
	}

	fmt.Printf("Issue atualizada com sucesso!\nURL: %v\n", result["html_url"])
	return nil
}
