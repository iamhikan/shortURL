package filestorage

import (
	"fmt"
	"log"
	"sync"
)

// Логика для хранения в файле
type FileStorage struct {
	fs               *FS
	IncrementStorage int
	mu               sync.RWMutex
}

func NewFileStorage(filePath string) *FileStorage {
	fs := NewFS(filePath)

	CurIncrement, err := fs.FindLastID()
	if err != nil {
		log.Fatalf("error reading last index: %s", err)
	}
	fmt.Printf("CurIncrement = %d\n", CurIncrement)

	return &FileStorage{
		fs:               fs,
		IncrementStorage: CurIncrement,
	}
}

func (fs *FileStorage) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.fs.file.Close()
}

func (fs *FileStorage) Get(id int) (string, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	link, found, err := fs.fs.FindLinkByID(id)
	if err != nil {
		log.Fatalf("error finding link by id: %s", err)
	}
	if !found {
		return "", false
	}
	return link, true
}

func (fs *FileStorage) Set(link string) int {
	fs.mu.Lock()
	defer func() {
		fs.IncrementStorage++
		fs.mu.Unlock()
	}()
	err := fs.fs.WriteData(&LinkData{
		OriginalLink: link,
		ID:           fs.IncrementStorage,
	})
	if err != nil {
		log.Fatalf("error writing to file: %s", err)
	}
	return fs.IncrementStorage
}
