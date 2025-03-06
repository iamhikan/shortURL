package inmemory

import "sync"

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
