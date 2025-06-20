package main

import (
	"ecommerce-stock-api/config"
	"ecommerce-stock-api/internal/order"
	"ecommerce-stock-api/internal/product"
	"ecommerce-stock-api/internal/shop"
	"ecommerce-stock-api/internal/user"
	"ecommerce-stock-api/internal/warehouse"
	"ecommerce-stock-api/jobs"
	"ecommerce-stock-api/pkg/db"
	"ecommerce-stock-api/pkg/logger"
	"ecommerce-stock-api/pkg/monitor"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	monitor.Init()

	// DB init
	database := db.ConnectPostgres()

	// User
	userRepo := user.NewRepository(database)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Product
	productRepo := product.NewRepository(database)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Shop
	shopRepo := shop.NewRepository(database)
	shopService := shop.NewService(shopRepo)
	shopHandler := shop.NewHandler(shopService)

	// Warehouse
	warehouseRepo := warehouse.NewRepository(database)
	warehouseService := warehouse.NewService(warehouseRepo)
	warehouseHandler := warehouse.NewHandler(warehouseService)

	// Order
	orderRepo := order.NewRepository(database)
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	r := gin.Default()

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(monitor.Handler()))

	// Optional middleware to track every request
	r.Use(func(c *gin.Context) {
		c.Next()
		monitor.HttpRequestsTotal.WithLabelValues(c.FullPath(), c.Request.Method).Inc()
	})

	// User
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Product
	r.GET("/products", productHandler.ListProducts)

	// Shop
	r.POST("/shops", shopHandler.Create)
	r.GET("/shops", shopHandler.List)

	// Warehouse
	r.POST("/warehouses", warehouseHandler.Create)
	r.POST("/warehouses/status", warehouseHandler.UpdateStatus)
	r.GET("/warehouses", warehouseHandler.List)
	r.POST("/warehouses/transfer", warehouseHandler.TransferStock)

	// Order
	r.POST("/order/checkout", orderHandler.Checkout)

	// Run server
	r.Run(":8080")

	jobs.StartStockReleaseWorker(database, 1*time.Minute)
}
