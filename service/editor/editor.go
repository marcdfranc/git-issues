package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"git-issues/domain"
)

type Editor interface {
	GetIssueContentFromEditor(initialTitle, initialBody string) (domain.Issue, error)
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

func (s *Service) GetIssueContentFromEditor(initialTitle, initialBody string) (domain.Issue, error) {
	tempFile, err := os.CreateTemp("", "ghissue-*.md")
	if err != nil {
		return domain.Issue{}, fmt.Errorf("could not create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := fmt.Sprintf("%s\n\n%s", initialTitle, initialBody)
	_, err = tempFile.WriteString(content)
	if err != nil {
		return domain.Issue{}, fmt.Errorf("could not write temp file: %v", err)
	}
	tempFile.Close()

	editor := s.getEditor()
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return domain.Issue{}, fmt.Errorf("error on exec text editor: %v", err)
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return domain.Issue{}, fmt.Errorf("error on read file: %v", err)
	}

	parts := strings.SplitN(string(editedContent), "\n", 2)
	response := domain.Issue{
		Title: strings.TrimSpace(parts[0]),
	}
	if len(parts) > 1 {
		response.Body = strings.TrimSpace(parts[1])
	}

	return response, nil
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
