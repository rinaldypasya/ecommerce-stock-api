package product

import (
	"gorm.io/gorm"
)

type Repository interface {
	ListProductsWithAvailableStock() ([]ProductWithStock, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

type ProductWithStock struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	TotalStock  int
}

func (r *repo) ListProductsWithAvailableStock() ([]ProductWithStock, error) {
	var result []ProductWithStock

	err := r.db.
		Table("products").
		Select("products.id, products.name, products.description, products.price, SUM(product_stocks.quantity - product_stocks.reserved_qty) as total_stock").
		Joins("LEFT JOIN product_stocks ON product_stocks.product_id = products.id").
		Group("products.id").
		Scan(&result).Error

	return result, err
}
