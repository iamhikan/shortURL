package main

import (
	"net/http"
	"short_url/internal/router"
	"short_url/internal/service"
)

func main() {
	r := router.SetupRouter()
	srv := service.New()
	router.Routes(r, srv)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
