package repository

import "module3/internal/entities"

type UrlRepository interface {
	Save(url, alias string) (*entities.Url, error)
	Update(url *entities.Url) (*entities.Url, error)
	DeleteById(id int64) error
	FindById(id int64) (*entities.Url, error)
	FindUrlStringByAlias(alias string) (string, error)
}
