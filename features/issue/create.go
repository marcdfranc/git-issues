package issue

import (
	"encoding/json"
	"fmt"
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


}