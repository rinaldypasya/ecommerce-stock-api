package jobs

import (
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/monitor"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartStockReleaseWorker(db *gorm.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			releaseExpiredOrders(db)
		}
	}()
}

func releaseExpiredOrders(db *gorm.DB) {
	var expiredOrders []models.Order
	now := time.Now()

	// Find expired PENDING orders
	if err := db.Where("status = ? AND expires_at <= ?", "PENDING", now).Find(&expiredOrders).Error; err != nil {
		log.Println("failed to find expired orders:", err)
		return
	}

	for _, order := range expiredOrders {
		err := db.Transaction(func(tx *gorm.DB) error {
			var items []models.OrderItem
			if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
				return err
			}

			for _, item := range items {
				var stock models.ProductStock
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					Where("product_id = ?", item.ProductID).
					First(&stock).Error; err != nil {
					return err
				}
				stock.ReservedQty -= item.Quantity
				if stock.ReservedQty < 0 {
					stock.ReservedQty = 0
				}
				if err := tx.Save(&stock).Error; err != nil {
					return err
				}

				monitor.StockReleaseTotal.Add(float64(item.Quantity))
			}

			// Cancel order
			if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("status", "CANCELLED").Error; err != nil {
				return err
			}
			log.Printf("Cancelled expired order ID %d and released stock\n", order.ID)
			return nil
		})

		if err != nil {
			log.Println("error processing expired order:", err)
		}
	}
}
