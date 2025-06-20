package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex"`
	Phone        string `gorm:"uniqueIndex"`
	PasswordHash string
	CreatedAt    time.Time
}

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
}

type ProductStock struct {
	ID          uint `gorm:"primaryKey"`
	ProductID   uint
	Product     Product `gorm:"foreignKey:ProductID"`
	WarehouseID uint
	Quantity    int
	ReservedQty int
}

type Shop struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Warehouse struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	ShopID   uint
	IsActive bool
}

type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Status    string // PENDING, PAID, CANCELLED
	CreatedAt time.Time
	ExpiresAt time.Time
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
	Quantity  int
}
