package stubs

import "git-issues/domain"

type EditorStub struct {
	GetIssueContentFromEditorFunc func(issue *domain.Issue) error
}

func (s *EditorStub) GetIssueContentFromEditor(issue *domain.Issue) error {
	if s.GetIssueContentFromEditorFunc != nil {
		return s.GetIssueContentFromEditorFunc(issue)
	}
	return nil
}
