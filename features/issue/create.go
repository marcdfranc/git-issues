package issue

import (
	"encoding/json"
	"fmt"

	"git-issues/client"
	"git-issues/domain"
)

func Create(config *domain.Config) {
	title, body, err := getIssueContentFromEditor(config, "", "")
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

	url := fmt.Sprintf("%s/repos/%s/%s/issues", config.APIBaseURL, config.Owner, config.Repo)
	response, err := client.MakeGitHubRequest(config, "POST", url, issue)
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
