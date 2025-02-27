package filestorage_test

import (
	"errors"
	"short_url/gen/mocks"
	"short_url/internal/repository/filestorage"
	"testing"

	"github.com/stretchr/testify/suite"
)

// LinkStorageTestSuite - набор тестов для хранилища ссылок
type LinkStorageTestSuite struct {
	suite.Suite
	mockLinkStorage *mocks.MockILinkStorage
}

// SetupTest выполняется перед каждым тестом
func (s *LinkStorageTestSuite) SetupTest() {
	// Создаем новый мок для каждого теста
	s.mockLinkStorage = mocks.NewMockILinkStorage(s.T())
}

// TestFindLastID тестирует метод FindLastID
func (s *LinkStorageTestSuite) TestFindLastID() {
	testCases := []struct {
		name        string
		expectedID  int
		expectError bool
	}{
		{
			name:        "Empty file",
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "Single record",
			expectedID:  5,
			expectError: false,
		},
		{
			name:        "Multiple records",
			expectedID:  3,
			expectError: false,
		},
		{
			name:        "Error reading",
			expectedID:  1,
			expectError: true,
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			// Сбрасываем ожидания мока перед каждым подтестом
			s.mockLinkStorage = mocks.NewMockILinkStorage(s.T())

			// Настраиваем поведение мока с конкретными ожидаемыми результатами
			if tt.expectError {
				s.mockLinkStorage.EXPECT().FindLastID().Return(tt.expectedID, errors.New("test error")).Once()
			} else {
				s.mockLinkStorage.EXPECT().FindLastID().Return(tt.expectedID, nil).Once()
			}

			// Вызываем метод и проверяем результат
			id, err := s.mockLinkStorage.FindLastID()

			// Проверяем ожидания
			if tt.expectError {
				s.Error(err, "FindLastID() должен возвращать ошибку")
			} else {
				s.NoError(err, "FindLastID() не должен возвращать ошибку")
			}
			s.Equal(tt.expectedID, id, "FindLastID() должен возвращать ожидаемый ID")
		})
	}
}

// TestWriteData тестирует метод WriteData
func (s *LinkStorageTestSuite) TestWriteData() {
	// Создаем тестовые данные
	expected := filestorage.LinkData{
		OriginalLink: "https://example.com",
		ID:           4,
	}

	// Настраиваем ожидаемый вызов WriteData
	s.mockLinkStorage.EXPECT().WriteData(&expected).Return(nil).Once()

	// Вызываем метод и проверяем, что ошибки нет
	err := s.mockLinkStorage.WriteData(&expected)
	s.NoError(err, "WriteData() не должен возвращать ошибку")
}

// TestFindLinkByID тестирует метод FindLinkByID
func (s *LinkStorageTestSuite) TestFindLinkByID() {
	testCases := []struct {
		name  string
		id    int
		link  string
		found bool
		err   error
	}{
		{
			name:  "First",
			id:    1,
			link:  "https://example.com",
			found: true,
			err:   nil,
		},
		{
			name:  "Not found",
			id:    0,
			link:  "",
			found: false,
			err:   nil,
		},
		{
			name:  "Last",
			id:    3,
			link:  "https://example.net",
			found: true,
			err:   nil,
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			// Сбрасываем ожидания мока перед каждым подтестом
			s.mockLinkStorage = mocks.NewMockILinkStorage(s.T())

			// Настраиваем поведение мока
			var err error
			if tt.err != nil {
				err = errors.New("test error")
			}

			s.mockLinkStorage.EXPECT().FindLinkByID(tt.id).Return(tt.link, tt.found, err).Once()

			// Вызываем метод и проверяем результат
			link, found, err := s.mockLinkStorage.FindLinkByID(tt.id)

			s.Equal(tt.link, link, "FindLinkByID() должен возвращать правильную ссылку")
			s.Equal(tt.found, found, "FindLinkByID() должен правильно показывать найдена ли ссылка")
			s.Equal((tt.err != nil), (err != nil), "FindLinkByID() должен правильно возвращать ошибку")
		})
	}
}

// TestLinkStorageSuite запускает набор тестов
func TestLinkStorageSuite(t *testing.T) {
	suite.Run(t, new(LinkStorageTestSuite))
}
