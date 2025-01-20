package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// LinkData предcтавляет данные ссылки
type LinkData struct {
	OriginalLink string `json:"original_link"`
	ID           int    `json:"ID"`
}

func FindLastID(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	if lastLine == "" {
		return 1, nil
	}
	var link LinkData
	if err := json.Unmarshal([]byte(lastLine), &link); err != nil {
		return 1, fmt.Errorf("ошибка разбора JSON: %w", err)
	}
	return link.ID, nil
}

// WriteData записывает данные в файл
func WriteData(file *os.File, ld *LinkData) error {
	data, err := json.Marshal(ld)
	if err != nil {
		return fmt.Errorf("ошибка при маршаллинге JSON %w", err)
	}
	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("ошибка при записи JSON в файл %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка при записи в файл символа переноса строки %w", err)
	}

	return writer.Flush()
}

func FindLinkByID(file *os.File, id int) (string, bool, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", false, fmt.Errorf("ошибка сброса указателя файла: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var link LinkData
		if err := json.Unmarshal([]byte(line), &link); err != nil {
			return "", false, fmt.Errorf("ошибка разбора JSON: %w", err)
		}
		if id == link.ID {
			return link.OriginalLink, true, nil
		}
	}
	if err = scanner.Err(); err != nil {
		return "", false, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return "", false, nil
}
