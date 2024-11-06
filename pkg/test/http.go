package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// doRequest отправляет HTTP-запрос и возвращает тестовый ответ.
func DoRequest(t *testing.T, handler http.HandlerFunc, method, url string, body []byte, setup func(req *http.Request) *http.Request) *httptest.ResponseRecorder {
	t.Helper() // помечаем как вспомогательную функцию для улучшения вывода ошибок в тестах

	var reqBody *bytes.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if setup != nil {
		req = setup(req)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}
