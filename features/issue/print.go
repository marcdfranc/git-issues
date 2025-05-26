package issue

import (
	"fmt"
	"io"

	"git-issues/domain"
)

const (
	strIssueFormat = "#%v - %s (%s)\n"
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
