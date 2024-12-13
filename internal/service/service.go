package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"short_url/config"
	"short_url/internal/repository"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Storage repository.IStorage
	Config  config.Config
}

func New(stor repository.IStorage) *Service {
	var Cfg config.Config
	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}
	return &Service{
		Storage: stor,
		Config:  Cfg,
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

	res := fmt.Sprintf("%s/%d", s.Config.BaseURL, id)
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

func (s *Service) CreateShortURLFromJSON(w http.ResponseWriter, r *http.Request) {
	var curURL CreateShortURLFromJSONReq
	if err := json.NewDecoder(r.Body).Decode(&curURL); err != nil {
		http.Error(w, "Некорректный формат", http.StatusBadRequest)
	}

	id := s.Storage.Set(curURL.URL)
	newShortLink := CreateShortURLFromJSONRes{
		Result: fmt.Sprintf("%s/%d", s.Config.BaseURL, id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newShortLink); err != nil {
		http.Error(w, "Ошибка при маршалинге ответа", http.StatusInternalServerError)
	}
}
