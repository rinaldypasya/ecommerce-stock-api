package warehouse

import (
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/monitor"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(warehouse *models.Warehouse) error
	UpdateStatus(id uint, active bool) error
	List() ([]models.Warehouse, error)
	TransferStock(input TransferInput) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db}
}

func (r *repo) Create(wh *models.Warehouse) error {
	return r.db.Create(wh).Error
}

func (r *repo) UpdateStatus(id uint, active bool) error {
	return r.db.Model(&models.Warehouse{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *repo) List() ([]models.Warehouse, error) {
	var wh []models.Warehouse
	err := r.db.Find(&wh).Error
	return wh, err
}

type TransferInput struct {
	ProductID       uint
	FromWarehouseID uint
	ToWarehouseID   uint
	Quantity        int
}

func (r *repo) TransferStock(input TransferInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if input.FromWarehouseID == input.ToWarehouseID {
			return errors.New("cannot transfer to the same warehouse")
		}

		// Lock and fetch source
		var from models.ProductStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", input.ProductID, input.FromWarehouseID).
			First(&from).Error; err != nil {
			return fmt.Errorf("source warehouse stock not found or lock error: %w", err)
		}

		if (from.Quantity - from.ReservedQty) < input.Quantity {
			return errors.New("insufficient stock in source warehouse")
		}

		from.Quantity -= input.Quantity
		if err := tx.Save(&from).Error; err != nil {
			return err
		}

		// Lock or create destination stock row
		var to models.ProductStock
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", input.ProductID, input.ToWarehouseID).
			First(&to).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			to = models.ProductStock{
				ProductID:   input.ProductID,
				WarehouseID: input.ToWarehouseID,
				Quantity:    input.Quantity,
				ReservedQty: 0,
			}
			return tx.Create(&to).Error
		} else if err != nil {
			return err
		}

		to.Quantity += input.Quantity

		monitor.StockTransfersTotal.Add(float64(input.Quantity))

		return tx.Save(&to).Error
	})
}
