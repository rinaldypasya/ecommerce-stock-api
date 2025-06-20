package shop

import (
	"ecommerce-stock-api/models"
)

type Service interface {
	Create(name string) error
	List() ([]models.Shop, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(name string) error {
	err := s.repo.Create(&models.Shop{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) List() ([]models.Shop, error) {
	shops, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return shops, nil
}
