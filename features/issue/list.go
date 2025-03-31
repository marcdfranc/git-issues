package issue

import (
	"encoding/json"
	"fmt"

	"git-issues/domain"
	"git-issues/service/client"
)

func List(config *domain.Config) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", config.APIBaseURL, config.Owner, config.Repo)
	response, err := client.MakeGitHubRequest(config, "GET", url, nil)
	if err != nil {
		fmt.Printf("request error: %v\n", err)
		return
	}

	var issues []map[string]interface{}
	err = json.Unmarshal(response, &issues)
	if err != nil {
		fmt.Printf("read body error: %v\n", err)
		return
	}

	if len(issues) == 0 {
		fmt.Println("no issues found.")
		return
	}

	fmt.Println("\nIssues:")
	for _, issue := range issues {
		fmt.Printf("#%v - %s (%s)\n", issue["number"], issue["title"], issue["state"])
	}
}
