package order_test

import (
	"ecommerce-stock-api/internal/order"
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/db"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentCheckout(t *testing.T) {
	database := db.ConnectPostgres()

	// Reset DB for clean test (use with caution!)
	_ = database.Exec("TRUNCATE orders, order_items, product_stocks, products RESTART IDENTITY CASCADE")

	// Setup stock: product 1 has only 1 unit
	product := models.Product{Name: "iPhone", Price: 999}
	database.Create(&product)

	stock := models.ProductStock{
		ProductID:   product.ID,
		WarehouseID: 1,
		Quantity:    1,
	}
	database.Create(&stock)

	// Setup order service
	orderRepo := order.NewRepository(database)
	orderService := order.NewService(orderRepo)

	// Concurrency test
	successCount := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			err := orderService.Checkout(userID, []models.OrderItem{
				{ProductID: product.ID, Quantity: 1},
			})
			mu.Lock()
			if err == nil {
				successCount++
				log.Printf("User %d: SUCCESS", userID)
			} else {
				log.Printf("User %d: FAILED - %v", userID, err)
			}
			mu.Unlock()
		}(uint(i + 1))
	}

	wg.Wait()

	assert.Equal(t, 1, successCount, "Only one checkout should succeed due to limited stock")
}
