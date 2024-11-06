package service

import (
	"fmt"
	"io"
	"net/http"
	"short_url/internal/repository"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	Storage repository.IStorage
}

func New(stor repository.IStorage) *Service {
	return &Service{
		Storage: stor,
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

	id := s.Storage.Set(string(body))

	res := fmt.Sprintf("%s/%d", "localhost:8080", id)
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

	originalURL, found := s.Storage.Get(id)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("incorrect link"))
		return
	}
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
