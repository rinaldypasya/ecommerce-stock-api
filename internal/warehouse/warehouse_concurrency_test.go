package warehouse_test

import (
	"ecommerce-stock-api/internal/warehouse"
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/db"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentWarehouseTransfers(t *testing.T) {
	database := db.ConnectPostgres()
	_ = database.Exec("TRUNCATE product_stocks, products RESTART IDENTITY CASCADE")

	// Seed product
	product := models.Product{Name: "Hard Drive"}
	database.Create(&product)

	// Initial stock: 100 in warehouse 1
	database.Create(&models.ProductStock{
		ProductID:   product.ID,
		WarehouseID: 1,
		Quantity:    100,
		ReservedQty: 0,
	})
	database.Create(&models.ProductStock{
		ProductID:   product.ID,
		WarehouseID: 2,
		Quantity:    0,
	})

	repo := warehouse.NewRepository(database)
	service := warehouse.NewService(repo)

	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0
	failures := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := service.TransferStock(product.ID, 1, 2, 15)
			mu.Lock()
			if err != nil {
				failures++
				log.Printf("Worker %d: FAILED (%v)", id, err)
			} else {
				success++
				log.Printf("Worker %d: SUCCESS", id)
			}
			mu.Unlock()
		}(i + 1)
	}

	wg.Wait()

	// Verify
	var source models.ProductStock
	var target models.ProductStock
	_ = database.Where("product_id = ? AND warehouse_id = 1", product.ID).First(&source)
	_ = database.Where("product_id = ? AND warehouse_id = 2", product.ID).First(&target)

	totalTransferred := target.Quantity
	assert.LessOrEqual(t, totalTransferred, 100)
	assert.Equal(t, 100, source.Quantity+target.Quantity)
	assert.Equal(t, 100-totalTransferred, source.Quantity)

	log.Printf("✅ Total success: %d, failures: %d, transferred: %d", success, failures, totalTransferred)
}
