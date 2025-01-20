package main

import (
	"log"
	"net/http"
	"short_url/config"
	"short_url/internal/repository"
	"short_url/internal/router"
	"short_url/internal/service"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки .env файла: %v", err)
	}

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("ошибка при парсинге конфигурации: %v", err)
	}

	var stor repository.IStorage
	if cfg.FileStoragePath == "" {
		stor = repository.New()
	} else {
		stor = repository.NewFileStorage(cfg.FileStoragePath)
		defer func() {
			if err := stor.Close(); err != nil {
				log.Fatalf("ошибка при закрытии файла: %v", err)
			}
		}()
	}

	srv := service.New(stor, cfg)
	r := router.SetupRouter()
	router.Routes(r, srv)

	err = http.ListenAndServe(srv.Config.ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
