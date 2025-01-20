package repository_test

import (
	"bufio"
	"encoding/json"
	"os"
	"short_url/internal/repository"
	"strings"
	"testing"
)

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()

	tempFile, err := os.CreateTemp("", "testfile-*")
	if err != nil {
		t.Fatalf("не удалось создать временный файл: %v", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		t.Fatalf("Не удалось записать в временный файл: %v", err)
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		t.Fatalf("Не удалось сбросить указатель: %v", err)
	}
	return tempFile
}

func TestFindLastID(t *testing.T) {

	for _, tt := range TestsForFindLast {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := createTempFile(t, tt.fileContent)
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()

			id, err := repository.FindLastID(tempFile)
			if (err != nil) != tt.expectError {
				t.Errorf("FindLastID() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if id != tt.expectedID {
				t.Errorf("FindLastID() = %v, expected %v", id, tt.expectedID)
			}

		})
	}
}

func TestWriteData(t *testing.T) {

	tempFile := createTempFile(t, "")
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	expected := repository.LinkData{
		OriginalLink: "https://example.com",
		ID:           4,
	}

	if err := repository.WriteData(tempFile, &expected); err != nil {
		t.Fatalf("ошибка при выполнении WriteData: %v", err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Не удалось сбросить указатель файла: %v", err)
	}

	reader := bufio.NewReader(tempFile)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Ошибка при чтении строки из файла: %v", err)
	}
	line = strings.TrimSpace(line)

	var result repository.LinkData
	if err := json.Unmarshal([]byte(line), &result); err != nil {
		t.Fatalf("Не удалось распарсить JSON: %v", err)
	}

	if result.OriginalLink != expected.OriginalLink || result.ID != expected.ID {
		t.Errorf("некорректные данные. Ожидались: link - %s, ID - %d. Получены: link - %s, ID - %d", expected.OriginalLink, expected.ID, result.OriginalLink, result.ID)
	}
}

func TestFindLinkByID(t *testing.T) {

	tempFile := createTempFile(t, `{"original_link":"https://example.com","ID":1}
{"original_link":"https://example.org","ID":2}
{"original_link":"https://example.net","ID":3}
`)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	for _, tt := range TestsForFindLinkByID {
		t.Run(tt.Name, func(t *testing.T) {
			if _, err := tempFile.Seek(0, 0); err != nil {
				t.Fatalf("не удалось сбросить указатель файла: %v", err)
			}
			link, found, err := repository.FindLinkByID(tempFile, tt.ID)
			if link != tt.link || found != tt.found || err != tt.err {
				t.Errorf("некорректные данные. Ожидались: found - %t, link - %s, error - %v. Получены:  found - %t, link - %s, error - %v", tt.found, tt.link, tt.err, found, link, err)
			}
		})

	}

}

var TestsForFindLast = []struct {
	name        string
	fileContent string
	expectedID  int
	expectError bool
}{
	{
		name:        "Empty file",
		fileContent: "",
		expectedID:  1,
		expectError: false,
	},
	{
		name:        "Single record",
		fileContent: `{"original_link":"https://example.com","ID":5}`,
		expectedID:  5,
		expectError: false,
	},
	{
		name: "Multiple records",
		fileContent: `{"original_link":"https://example.com","ID":1}
			{"original_link":"https://example.org","ID":2}
			{"original_link":"https://example.net","ID":3}
			`,
		expectedID:  3,
		expectError: false,
	},
	{
		name: "Invalid JSON",
		fileContent: `{"original_link":"https://example.com","ID":1}
			invalid json
			`,
		expectedID:  1, // По логике, возвращается 1 при ошибке разбора последней строки
		expectError: true,
	},
	{
		name: "Last line empty",
		fileContent: `{"original_link":"https://example.com","ID":1}
			`,
		expectedID:  1,
		expectError: false,
	},
}

var TestsForFindLinkByID = []struct {
	Name  string
	ID    int
	link  string
	found bool
	err   error
}{
	{
		Name:  "First",
		ID:    1,
		link:  "https://example.com",
		found: true,
		err:   nil,
	},
	{
		Name:  "Not found",
		ID:    0,
		link:  "",
		found: false,
		err:   nil,
	},
	{
		Name:  "Last",
		ID:    3,
		link:  "https://example.net",
		found: true,
		err:   nil,
	},
}
