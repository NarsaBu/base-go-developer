package mocks

import (
	"module3/internal/entities"

	"github.com/stretchr/testify/mock"
)

type MockUrlRepository struct {
	mock.Mock
}

func (m *MockUrlRepository) Save(url, alias string) (*entities.Url, error) {
	args := m.Called(url, alias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Url), args.Error(1)
}

func (m *MockUrlRepository) Update(url *entities.Url) (*entities.Url, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Url), args.Error(1)
}

func (m *MockUrlRepository) DeleteById(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUrlRepository) FindById(id int64) (*entities.Url, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Url), args.Error(1)
}

func (m *MockUrlRepository) FindUrlStringByAlias(alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}
