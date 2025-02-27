package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LinkData предcтавляет данные ссылки
type LinkData struct {
	OriginalLink string `json:"original_link"`
	ID           int    `json:"ID"`
}

// ILinkStorage определяет интерфейс для хранилища ссылок
type ILinkStorage interface {
	// FindLastID находит последний использованный ID
	FindLastID() (int, error)
	// WriteData записывает данные ссылки в хранилище
	WriteData(ld *LinkData) error
	// FindLinkByID ищет ссылку по ID
	FindLinkByID(id int) (string, bool, error)
	// Close закрывает хранилище
	Close() error
}

// FS реализует ILinkStorage с использованием файловой системы
type FS struct {
	file *os.File
}

// NewFS создаёт новый экземпляр файлового хранилища
func NewFS(filePath string) *FS {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}

	return &FS{file: file}
}

// FindLastID находит последний использованный ID
func (f *FS) FindLastID() (int, error) {
	scanner := bufio.NewScanner(f.file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}
	if lastLine == "" {
		return 1, nil
	}
	var link LinkData
	if err := json.Unmarshal([]byte(lastLine), &link); err != nil {
		return 1, fmt.Errorf("error parsing JSON: %w", err)
	}
	return link.ID, nil
}

// WriteData записывает данные в файл
func (f *FS) WriteData(ld *LinkData) error {
	data, err := json.Marshal(ld)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	writer := bufio.NewWriter(f.file)
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error writing newline to file: %w", err)
	}

	return writer.Flush()
}

// FindLinkByID ищет ссылку по ID
func (f *FS) FindLinkByID(id int) (string, bool, error) {
	_, err := f.file.Seek(0, 0)
	if err != nil {
		return "", false, fmt.Errorf("error resetting file pointer: %w", err)
	}

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		line := scanner.Text()
		var link LinkData
		if err := json.Unmarshal([]byte(line), &link); err != nil {
			return "", false, fmt.Errorf("error parsing JSON: %w", err)
		}
		if id == link.ID {
			return link.OriginalLink, true, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return "", false, fmt.Errorf("error reading file: %w", err)
	}

	return "", false, nil
}

// Close закрывает файл
func (f *FS) Close() error {
	return f.file.Close()
}
