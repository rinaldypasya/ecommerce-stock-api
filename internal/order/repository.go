package order

import (
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/monitor"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateOrder(order *models.Order, items []models.OrderItem) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) CreateOrder(order *models.Order, items []models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range items {
			var stock models.ProductStock
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("product_id = ? AND (quantity - reserved_qty) >= ?", item.ProductID, item.Quantity).
				First(&stock).Error
			if err != nil {
				return err
			}
			stock.ReservedQty += item.Quantity
			if err := tx.Save(&stock).Error; err != nil {
				return err
			}
			monitor.StockReservedTotal.Add(float64(item.Quantity))

			item.OrderID = order.ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
