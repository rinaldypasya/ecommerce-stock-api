package order_test

import (
	"ecommerce-stock-api/internal/order"
	"ecommerce-stock-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	lastOrder  *models.Order
	lastItems  []models.OrderItem
	failCreate bool
}

func (m *mockRepo) CreateOrder(o *models.Order, items []models.OrderItem) error {
	if m.failCreate {
		return assert.AnError
	}
	m.lastOrder = o
	m.lastItems = items
	return nil
}

func TestCheckout_Success(t *testing.T) {
	mock := &mockRepo{}
	svc := order.NewService(mock)

	items := []models.OrderItem{
		{ProductID: 1, Quantity: 2},
	}

	err := svc.Checkout(42, items)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), mock.lastOrder.UserID)
	assert.Len(t, mock.lastItems, 1)
	assert.Equal(t, 2, mock.lastItems[0].Quantity)
}
