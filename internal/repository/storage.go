package repository

import "sync"

type Storage struct {
	LinksStorage     map[int]string
	IncrementStorage int // последний айди
	mu               sync.RWMutex
}

func New() *Storage {
	return &Storage{
		LinksStorage:     make(map[int]string),
		IncrementStorage: 1,
	}
}

func (s *Storage) Get(id int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, found := s.LinksStorage[id]
	return v, found
}

func (s *Storage) Set(link string) int {
	s.mu.Lock()
	defer func() {
		s.IncrementStorage++
		s.mu.Unlock()
	}()
	s.LinksStorage[s.IncrementStorage] = link
	return s.IncrementStorage
}
