package usecase

import (
	"module3/internal/dto"
	"module3/internal/entities"
	"module3/internal/repository"
)

type UrlService struct {
	repository repository.UrlRepository
}

func NewUrlService(urlRepository repository.UrlRepository) *UrlService {
	return &UrlService{repository: urlRepository}
}

func (us *UrlService) CreateUrl(url, alias string) (*dto.UrlResponse, error) {
	createdUser, err := us.repository.Save(url, alias)
	if err != nil {
		return nil, err
	}

	return &dto.UrlResponse{
		Id:    createdUser.Id,
		Url:   createdUser.Url,
		Alias: createdUser.Alias,
	}, nil
}

func (us *UrlService) UpdateUrl(urlUpdateRequest *dto.UrlUpdateRequest) (*dto.UrlResponse, error) {
	urlEntity := entities.Url{
		Id:    urlUpdateRequest.Id,
		Url:   urlUpdateRequest.Url,
		Alias: urlUpdateRequest.Alias,
	}

	updatedUser, err := us.repository.Update(&urlEntity)
	if err != nil {
		return nil, err
	}

	return &dto.UrlResponse{
		Id:    updatedUser.Id,
		Url:   updatedUser.Url,
		Alias: updatedUser.Alias,
	}, nil
}

func (us *UrlService) DeleteById(id int64) error {
	err := us.repository.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UrlService) FindById(id int64) (*dto.UrlResponse, error) {
	foundUser, err := us.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	return &dto.UrlResponse{
		Id:    foundUser.Id,
		Url:   foundUser.Url,
		Alias: foundUser.Alias,
	}, nil
}

func (us *UrlService) FindUrlStringByAlias(alias string) (string, error) {
	urlToRedirect, err := us.repository.FindUrlStringByAlias(alias)
	if err != nil {
		return "", err
	}

	return urlToRedirect, nil
}
