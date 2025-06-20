package shop_test

import (
	"ecommerce-stock-api/internal/shop"
	"ecommerce-stock-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createdShops []models.Shop
	failCreate   bool
	returnShops  []models.Shop
	failList     bool
}

func (m *mockRepo) Create(shop *models.Shop) error {
	if m.failCreate {
		return assert.AnError
	}
	m.createdShops = append(m.createdShops, *shop)
	return nil
}

func (m *mockRepo) List() ([]models.Shop, error) {
	if m.failList {
		return nil, assert.AnError
	}
	return m.returnShops, nil
}

func TestCreateShop_Success(t *testing.T) {
	mock := &mockRepo{}
	svc := shop.NewService(mock)

	err := svc.Create("Toko ABC")
	assert.NoError(t, err)
	assert.Len(t, mock.createdShops, 1)
	assert.Equal(t, "Toko ABC", mock.createdShops[0].Name)
}

func TestListShops_Success(t *testing.T) {
	mock := &mockRepo{returnShops: []models.Shop{{ID: 1, Name: "Toko A"}}}
	svc := shop.NewService(mock)

	shops, err := svc.List()
	assert.NoError(t, err)
	assert.Len(t, shops, 1)
	assert.Equal(t, "Toko A", shops[0].Name)
}
