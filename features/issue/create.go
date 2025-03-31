package issue

import (
	"encoding/json"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
)

type Feature struct {
	config *domain.Config
	editor editor.Editor
}

func New(config *domain.Config, editor editor.Editor) *Feature {
	return &Feature{
		config: config,
		editor: editor,
	}
}

func (f *Feature) Create() {
	title, body, err := f.editor.GetIssueContentFromEditor("", "")
	if err != nil {
		fmt.Printf("could not edit issue: %v\n", err)
		return
	}

	if title == "" {
		fmt.Println("title is required.")
		return
	}

	issue := domain.Issue{
		Title: title,
		Body:  body,
	}

	url := fmt.Sprintf("%s/repos/%s/%s/issues", f.config.APIBaseURL, f.config.Owner, f.config.Repo)
	response, err := client.MakeGitHubRequest(f.config, "POST", url, issue)
	if err != nil {
		fmt.Printf("Could not create issue: %v\n", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Printf("error on process response: %v\n", err)
		return
	}

	fmt.Printf("Issue created with success!\nNumber: %v\nURL: %v\n", result["number"], result["html_url"])
}
