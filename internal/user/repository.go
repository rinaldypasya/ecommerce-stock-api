package user

import (
	"ecommerce-stock-api/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindByEmailOrPhone(identifier string) (*models.User, error)
	Create(user *models.User) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) FindByEmailOrPhone(identifier string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ? OR phone = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) Create(user *models.User) error {
	return r.db.Create(user).Error
}
