package issue

import (
	"errors"
	"os"
	"testing"

	"git-issues/domain"
	"git-issues/service/client"
	"git-issues/service/editor"
	"git-issues/testdata/stubs"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	// Arrange

	f := Feature{config: &domain.Config{}}

	tests := []struct {
		name       string
		clientStub client.GitHubClient
		editorStub editor.Editor
		want       error
	}{
		{
			name: "successful creation",
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(initialTitle, initialBody string) (domain.Issue, error) {
					return domain.Issue{
						Title: "Test Issue",
						Body:  "Test Body",
					}, nil
				},
			},
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data interface{}) ([]byte, error) {
					return []byte(`{"number":1,"html_url":"http://example.com/issue/1","title":"Test Issue","body":"Test Body"}`), nil
				},
			},
			want: nil,
		},
		{
			name: "empty title",
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(initialTitle, initialBody string) (domain.Issue, error) {
					return domain.Issue{
						Body: "Test Body",
					}, nil
				},
			},
			clientStub: &stubs.ClientStub{},
			want:       errTitleRequired,
		},
		{
			name: "editor error",
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(initialTitle, initialBody string) (domain.Issue, error) {
					return domain.Issue{}, domain.ErrEditor
				},
			},
			clientStub: &stubs.ClientStub{},
			want:       domain.ErrEditor,
		},
		{
			name: "api error",
			editorStub: &stubs.EditorStub{
				GetIssueContentFromEditorFunc: func(initialTitle, initialBody string) (domain.Issue, error) {
					return domain.Issue{
						Title: "Test Issue",
						Body:  "Test Body",
					}, nil
				},
			},
			clientStub: &stubs.ClientStub{
				MakeRequestFunc: func(method, url string, data interface{}) ([]byte, error) {
					return []byte{}, domain.ErrApi
				},
			},
			want: domain.ErrApi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f.editor = tt.editorStub
			f.client = tt.clientStub
			_, err := f.Create()
			if !errors.Is(err, tt.want) {
				t.Errorf("unxpected error got: %v\nbut want: %v", err, tt.want)
			}
		})
	}
}
