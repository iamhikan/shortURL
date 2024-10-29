package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type Service struct {
	LinksStorage     map[int]string
	IncrementStorage int // последний айди
	mu               sync.RWMutex
}

func New() *Service {
	return &Service{
		LinksStorage: make(map[int]string),
	}
}

//все ниже в этом файле - handlers
// http обработчики

func (s *Service) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

// ручка на POST /
func (s *Service) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	s.mu.Lock()
	s.LinksStorage[s.IncrementStorage] = string(body)
	defer func() {
		s.IncrementStorage++
		s.mu.Unlock()
	}()

	res := fmt.Sprintf("%s/%d", "localhost:8080", s.IncrementStorage)
	w.Write([]byte(res))
}

// ручка на GET /{id}
func (s *Service) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if _, found := s.LinksStorage[id]; !found {
		w.Write([]byte("incorrect link"))
	}
	w.WriteHeader(307)
	w.Header().Set("Location", s.LinksStorage[id])

}
