package router

import (
	"short_url/internal/service"

	"github.com/go-chi/chi/v5"
)

func Routes(router *chi.Mux, srv *service.Service) {
	// пример обработчика запроса, в его качестве выступает функция MainPage
	// обрати внимание, что мы передаем именно саму функцию как параметр
	router.Get("/", srv.MainPage)
	router.Post("/", srv.CreateShortURL)
	router.Get("/{id}", srv.GetOriginalURL)
}

// HTTP
// client -> server; server -> client

// HTTP Req
// Headers - заголовки. Служебная информация о запросе
// Body - тело запроса. JSON.... etc - любой набор байтов!

// Популярные заголовки
// Path - путь до страницы: /api/v1/page
// Method - GET, POST, PUT, PATCH, DELETE
// Date - дата время запроса
// Content-Type: mime-type. text/json, text/plain, text/javascript/, image/png