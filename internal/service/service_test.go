package service_test

import (
	"context"
	"fmt"
	"net/http"
	"short_url/gen/mocks"
	"short_url/internal/service"
	"short_url/pkg/test"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type serviceSuite struct {
	suite.Suite
	mockStorage *mocks.MockIStorage

	service *service.Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &serviceSuite{})
}

func (s *serviceSuite) SetupSuite() {
	s.mockStorage = mocks.NewMockIStorage(s.T())

	s.service = service.New(s.mockStorage)
}

func (s *serviceSuite) TestCreateShortURL() {
	reqBody := []byte("some_url")
	expectedID := 1

	s.mockStorage.EXPECT().Set(string(reqBody)).Return(expectedID)

	res := test.DoRequest(s.T(), s.service.CreateShortURL, http.MethodPost, "/", reqBody, nil)

	expectedRes := []byte(fmt.Sprintf("localhost:8080/%d", expectedID))

	s.Equal(http.StatusCreated, res.Code)
	s.Equal(expectedRes, res.Body.Bytes())
}

func (s *serviceSuite) TestGetOriginalURL() {
	reqBody := "http://some_url"
	expectedID := 1

	s.mockStorage.EXPECT().Get(expectedID).Return(reqBody, true)

	res := test.DoRequest(s.T(), s.service.GetOriginalURL, http.MethodGet, fmt.Sprintf("/%d", expectedID), nil, func(req *http.Request) *http.Request {
		// Создаем контекст маршрута и добавляем параметр `id`
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(expectedID))
		return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	})

	s.Equal(http.StatusTemporaryRedirect, res.Code)
	s.Equal(reqBody, res.Header().Get("Location"))
}