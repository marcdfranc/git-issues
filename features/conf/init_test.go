package conf

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"git-issues/domain"
	"git-issues/testdata/data"
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

	want := data.DefaultConfig

	ft := New()

	ft.writeFile = fakeWriteFile
	ft.reader = reader

	// ACT
	err := ft.Init()
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

func TestInitWriteFileError(t *testing.T) {
	// ARRANGE
	input := "tkn\nowner\nrepo\neditor\n"
	reader := strings.NewReader(input)

	f := New()
	f.reader = reader
	f.writeFile = func(filename string, data []byte, perm os.FileMode) error {
		return errors.New("disk full")
	}

	// ACT
	err := f.Init()

	// ASSERT
	if err == nil {
		t.Fatalf("expected error when writeFile fails, got nil")
	}
}

func TestGetConfigLoadsFile(t *testing.T) {
	// ARRANGE
	want := data.DefaultConfig
	data, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal want config: %v", err)
	}
	// write file that loadConfig will read
	if err := os.WriteFile(domain.ConfigFile, data, 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(domain.ConfigFile) })

	f := New()

	// ACT
	got, err := f.GetConfig()
	if err != nil {
		t.Fatalf("unexpected error from GetConfig: %v", err)
	}

	// ASSERT
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("loaded config does not match want. got=%v want=%v", *got, want)
	}
}

func TestGetConfigReturnsCached(t *testing.T) {
	// ARRANGE
	cached := &domain.Config{Token: "cachedToken"}
	f := New()
	f.config = cached

	// ACT
	got, err := f.GetConfig()
	if err != nil {
		t.Fatalf("unexpected error from GetConfig with cached config: %v", err)
	}

	// ASSERT: same pointer
	if got != cached {
		t.Fatalf("GetConfig did not return cached config pointer")
	}
}

func TestGetConfigReadFileError(t *testing.T) {
	// Ensure config file does not exist
	_ = os.Remove(domain.ConfigFile)
	t.Cleanup(func() { _ = os.Remove(domain.ConfigFile) })

	f := New()

	// ACT
	_, err := f.GetConfig()

	// ASSERT: should wrap errReadConfig
	if err == nil {
		t.Fatalf("expected error when config file is missing, got nil")
	}
	if !errors.Is(err, errReadConfig) {
		t.Fatalf("expected error to be wrapped with errReadConfig; got: %v", err)
	}
}

func TestGetConfigInvalidJSON(t *testing.T) {
	// ARRANGE: write invalid json
	if err := os.WriteFile(domain.ConfigFile, []byte("{ invalid json }"), 0600); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(domain.ConfigFile) })

	f := New()

	// ACT
	_, err := f.GetConfig()

	// ASSERT: should wrap errReadConfig due to unmarshal error
	if err == nil {
		t.Fatalf("expected error for invalid json, got nil")
	}
	if !errors.Is(err, errReadConfig) {
		t.Fatalf("expected error to be wrapped with errReadConfig for invalid json; got: %v", err)
	}
}
