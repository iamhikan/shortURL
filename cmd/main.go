package main

import (
	"net/http"
	"short_url/internal/repository"
	"short_url/internal/router"
	"short_url/internal/service"
)

func main() {
	r := router.SetupRouter()
	stor := repository.New()
	srv := service.New(stor)
	router.Routes(r, srv)
	err := http.ListenAndServe(srv.Config.ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
