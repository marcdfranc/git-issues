package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"git-issues/domain"
)

var (
	errCreateTmpFile = errors.New("could not create temp file")
	errWriteTempFile = errors.New("could not write temp file")
	errExecEditor    = errors.New("could not exec editor")
	errReadEditor    = errors.New("could not read editor output")
)

type Editor interface {
	GetIssueContentFromEditor(issue *domain.Issue) error
}

type Service struct {
	config *domain.Config
}

func New(config *domain.Config) *Service {
	return &Service{
		config: config,
	}
}

var goos = runtime.GOOS

func (s *Service) GetIssueContentFromEditor(issue *domain.Issue) error {
	tempFile, err := os.CreateTemp("", "ghissue-*.md")
	if err != nil {
		return errors.Join(errCreateTmpFile, err)
	}
	defer os.Remove(tempFile.Name())

	content := fmt.Sprintf("%s\n\n%s", issue.Title, issue.Body)
	_, err = tempFile.WriteString(content)
	if err != nil {
		return errors.Join(errWriteTempFile, err)
	}
	tempFile.Close()

	editor := s.getEditor()
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return errors.Join(errExecEditor, err)
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return errors.Join(errReadEditor, err)
	}

	parts := strings.SplitN(string(editedContent), "\n", 2)
	if issue.Title == "" {
		issue.Title = parts[0]
	}

	if len(parts) > 1 {
		issue.Body = strings.TrimSpace(parts[1])
	}

	return nil
}

func (s *Service) getEditor() string {
	if s.config.Editor != "" {
		return s.config.Editor
	}

	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// Default editors per OS
	if goos == "windows" {
		return "notepad"
	}
	return "vi"
}
