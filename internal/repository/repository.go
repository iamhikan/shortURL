package repository

import (
	"short_url/config"
	"short_url/internal/repository/filestorage"
	"short_url/internal/repository/inmemory"
)

// New создает и возвращает хранилище в зависимости от конфигурации
func New(cfg config.Config) IStorage {
	if cfg.FileStoragePath == "" {
		return inmemory.New()
	}
	return filestorage.NewFileStorage(cfg.FileStoragePath)
}
