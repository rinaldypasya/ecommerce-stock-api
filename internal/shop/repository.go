package shop

import (
	"ecommerce-stock-api/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(shop *models.Shop) error
	List() ([]models.Shop, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

func (r *repo) List() ([]models.Shop, error) {
	var shops []models.Shop
	err := r.db.Find(&shops).Error
	return shops, err
}
