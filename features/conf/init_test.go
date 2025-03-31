package conf

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"git-issues/domain"
)

func TestInitConfig(t *testing.T) {
	// ARRANGE
	// user inputs
	input := "myToken\nmyOwner\nmyRepo\nmyEditor\n"
	reader := strings.NewReader(input)

	// "Fake" FileWriter store data in a variable
	var writtenData []byte
	fakeWriteFile := func(filename string, data []byte, perm os.FileMode) error {
		if filename != domain.ConfigFile {
			t.Errorf("unexpected file name: want %s, got %s", domain.ConfigFile, filename)
		}
		writtenData = data
		return nil
	}

	want := domain.Config{
		Token:      "myToken",
		Owner:      "myOwner",
		Repo:       "myRepo",
		Editor:     "myEditor",
		APIBaseURL: apiBaseUrl,
	}

	ft := New()

	ft.writeFile = fakeWriteFile
	ft.reader = reader

	// ACT
	err := ft.InitConfig()
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Decode fields.
	var got domain.Config
	err = json.Unmarshal(writtenData, &got)
	if err != nil {
		t.Fatalf("error on json decoding: %v", err)
	}

	// ASSERT
	if want != got {
		t.Errorf("unexpected configuration")
	}
}
