package main

import (
	"fmt"
	"net/http"
	"short_url/internal/repository"
	"short_url/internal/router"
	"short_url/internal/service"
)

func main() {
	r := router.SetupRouter()
	stor := repository.New()
	srv := service.New(stor)
	fmt.Printf("BaseURL = %s, AddressServer = %s", srv.Config.BaseURL, srv.Config.ServerAddress)
	router.Routes(r, srv)
	err := http.ListenAndServe(srv.Config.ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
