package order

import (
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/monitor"
	"time"
)

type Service interface {
	Checkout(userID uint, items []models.OrderItem) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Checkout(userID uint, items []models.OrderItem) error {
	order := &models.Order{
		UserID:    userID,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	err := s.repo.CreateOrder(order, items)
	if err != nil {
		return err
	}

	monitor.OrdersCreatedTotal.Inc()

	return nil
}
