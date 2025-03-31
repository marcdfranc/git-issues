package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"git-issues/domain"
)

const apiBaseUrl = "https://api.github.com"

type Feature struct {
	reader    io.Reader
	writeFile func(filename string, data []byte, perm os.FileMode) error
}

func New() *Feature {
	return &Feature{
		writeFile: os.WriteFile,
		reader:    os.Stdin,
	}
}

func (f *Feature) InitConfig() error {
	reader := bufio.NewReader(f.reader)

	fmt.Print("GitHub Personal Access Token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	token = strings.TrimSpace(token)

	fmt.Print("Repository Owner (username/organization): ")
	owner, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	owner = strings.TrimSpace(owner)

	fmt.Print("Repository Name: ")
	repo, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	repo = strings.TrimSpace(repo)

	fmt.Print("Default text editor (empty for system default): ")
	editor, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	editor = strings.TrimSpace(editor)

	config := domain.Config{
		Token:      token,
		Owner:      owner,
		Repo:       repo,
		Editor:     editor,
		APIBaseURL: apiBaseUrl,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not generat conf: %w\n", err)
	}

	err = f.writeFile(domain.ConfigFile, configData, 0600)
	if err != nil {
		return fmt.Errorf("could not save conf: %w\n", err)
	}

	fmt.Println("conf created with success!")
	return nil
}
