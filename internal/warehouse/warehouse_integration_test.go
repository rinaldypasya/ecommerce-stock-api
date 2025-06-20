package warehouse_test

import (
	"ecommerce-stock-api/internal/warehouse"
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarehouseStockTransfer_Integration(t *testing.T) {
	database := db.ConnectPostgres()

	// Reset data
	database.Exec("TRUNCATE products, product_stocks RESTART IDENTITY CASCADE")

	// Create product
	product := models.Product{Name: "Laptop"}
	database.Create(&product)

	// Create stock in warehouse 1
	from := models.ProductStock{
		ProductID:   product.ID,
		WarehouseID: 1,
		Quantity:    10,
	}
	database.Create(&from)

	// Ensure warehouse 2 has no stock
	to := models.ProductStock{
		ProductID:   product.ID,
		WarehouseID: 2,
		Quantity:    0,
	}
	database.Create(&to)

	// Run transfer
	repo := warehouse.NewRepository(database)
	svc := warehouse.NewService(repo)
	err := svc.TransferStock(product.ID, 1, 2, 4)

	assert.NoError(t, err)

	var fromCheck, toCheck models.ProductStock
	database.Where("product_id = ? AND warehouse_id = ?", product.ID, 1).First(&fromCheck)
	database.Where("product_id = ? AND warehouse_id = ?", product.ID, 2).First(&toCheck)

	assert.Equal(t, 6, fromCheck.Quantity)
	assert.Equal(t, 4, toCheck.Quantity)
}
