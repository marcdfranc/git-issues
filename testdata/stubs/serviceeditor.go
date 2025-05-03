package stubs

import "git-issues/domain"

type EditorStub struct {
	GetIssueContentFromEditorFunc func(initialTitle, initialBody string) (domain.Issue, error)
}

func (s *EditorStub) GetIssueContentFromEditor(initialTitle, initialBody string) (domain.Issue, error) {
	if s.GetIssueContentFromEditorFunc != nil {
		return s.GetIssueContentFromEditorFunc(initialTitle, initialBody)
	}
	return domain.Issue{}, nil
}
