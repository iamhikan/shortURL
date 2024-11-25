package repository_test

import (
	"short_url/internal/repository"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	Storage *repository.Storage
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &StorageSuite{})
}

func (s *StorageSuite) SetupSuite() {
	storage := repository.New()

	s.Storage = storage
}

// любой тест пишется так:
// есть какие то входные данные
// есть те данные/поведение, которое мы ожидаем увидеть
// в качестве реакции на входные данные
// далее, запускаем функцию с входными данными
// получаем реакцию (например, выходные данные)
// оцениваем выходные данные с ожидаемыми
// если все совпало, то тест пройден
func (s *StorageSuite) TestGet() {
	// мы должны установить какие то существующие данные
	// попробовать их получить
	// если получили то, что ожидали - тест пройден
	expectedLink := "some_address"
	id := s.Storage.Set(expectedLink)

	currentLink, found := s.Storage.Get(id)
	s.True(found)
	s.Equal(expectedLink, currentLink)
}

func (s *StorageSuite) TestSet() {
	link := "some_adress"
	expectedId := 1

	currId := s.Storage.Set(link)
	s.Equal(expectedId, currId)

	// получим записанный адрес
	currLink, found := s.Storage.Get(currId)
	s.True(found)
	s.Equal(link, currLink)
}
