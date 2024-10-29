package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// функция создает роутер
// роутер - это структура и набор методов, которые
// обрабатывают входящий http запрос и в зависимости
// от метода и пути запроса вызывают указанную функцию для его обработки
func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	// указываем, какие middleware будем использовать
	// middleware - это функция, которая является прослойкой
	// и содержит некоторый полезный код, который выполняется
	// сквозным образом (до того момента, как будет вызван основной обработчик запроса)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	return router
}
