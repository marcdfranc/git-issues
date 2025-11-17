package issue

import (
	"fmt"
	"io"

	"git-issues/domain"
)

const (
	strIssueFormat       = "#%v - %s (%s)\n"
	strDetailIssueFormat = "\nIssue #%d\nTitle: %s\nState: %s\nBody:\n%s\n"
)

func PrintIssues(w io.Writer, issues []domain.Issue) error {
	_, err := fmt.Fprintln(w, "\nIssues:")
	if err != nil {
		return err
	}

	for _, i := range issues {
		_, err = fmt.Fprintf(w, strIssueFormat, i.Number, i.Title, i.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintIssue(w io.Writer, issue *domain.Issue) error {
	_, err := fmt.Fprintf(w, strDetailIssueFormat, issue.Number, issue.Title, issue.State, issue.Body)
	if err != nil {
		return err
	}
	return nil
}
