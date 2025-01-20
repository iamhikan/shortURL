package repository

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Логика для хранения данных в локальной сессии
type LocalStorage struct {
	LinksStorage     map[int]string
	IncrementStorage int // последний айди
	mu               sync.RWMutex
}

func New() *LocalStorage {
	return &LocalStorage{
		LinksStorage:     make(map[int]string),
		IncrementStorage: 1,
	}
}

func (s *LocalStorage) Get(id int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, found := s.LinksStorage[id]
	return v, found
}

func (s *LocalStorage) Set(link string) int {
	s.mu.Lock()
	defer func() {
		s.IncrementStorage++
		s.mu.Unlock()
	}()
	s.LinksStorage[s.IncrementStorage] = link
	return s.IncrementStorage
}

func (s *LocalStorage) Close() error {
	return nil
}

// Логика для хранения в файле
type FileStorage struct {
	file             *os.File
	IncrementStorage int
	mu               sync.RWMutex
}

func NewFileStorage(FileStoragePath string) *FileStorage {
	file, err := os.OpenFile(FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("ошибка при открытии файла: %s", err)
	}
	CurIncrement, err := FindLastID(file)
	if err != nil {
		log.Fatalf("ошибка при считывании последнего индекса: %s", err)
	}
	fmt.Printf("CurIncrement = %d\n", CurIncrement)

	return &FileStorage{
		file:             file,
		IncrementStorage: CurIncrement,
	}
}

func (fs *FileStorage) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.file.Close()
}

func (fs *FileStorage) Get(id int) (string, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	link, found, err := FindLinkByID(fs.file, id)
	if err != nil {
		log.Fatalf("ошибка при поиске из файла: %s", err)
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
	err := WriteData(fs.file, &LinkData{
		OriginalLink: link,
		ID:           fs.IncrementStorage,
	})
	if err != nil {
		log.Fatalf("ошибка записи в файл: %s", err)
	}
	return fs.IncrementStorage
}
