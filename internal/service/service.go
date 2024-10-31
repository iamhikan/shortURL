package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	Storage *Storage
}

type Storage struct {
	LinksStorage     map[int]string
	IncrementStorage int // последний айди
	mu               sync.RWMutex
}

func New() *Service {
	return &Service{
		Storage: &Storage{LinksStorage: make(map[int]string)},
	}
}

//все ниже в этом файле - handlers
// http обработчики

// ручка на POST /
func (s *Service) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	s.Storage.mu.Lock()
	s.Storage.LinksStorage[s.Storage.IncrementStorage] = string(body)
	defer func() {
		s.Storage.IncrementStorage++
		s.Storage.mu.Unlock()
	}()

	res := fmt.Sprintf("%s/%d", "localhost:8080", s.Storage.IncrementStorage)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

// ручка на GET /{id}
func (s *Service) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("incorrect format of id"))
		return
	}

	s.Storage.mu.RLock()
	defer func() {
		s.Storage.mu.RUnlock()
	}()
	OriginalURL, found := s.Storage.LinksStorage[id]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("incorrect link"))
		return
	}
	if !strings.HasPrefix(OriginalURL, "http://") && !strings.HasPrefix(OriginalURL, "https://") {
		OriginalURL = "http://" + OriginalURL
	}
	w.Header().Set("Location", OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
