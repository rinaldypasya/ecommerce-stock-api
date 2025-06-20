package product_test

import (
	"ecommerce-stock-api/internal/product"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	data []product.ProductWithStock
	err  error
}

func (m *mockRepo) ListProductsWithAvailableStock() ([]product.ProductWithStock, error) {
	return m.data, m.err
}

func TestGetProductList_Success(t *testing.T) {
	mock := &mockRepo{
		data: []product.ProductWithStock{
			{ID: 1, Name: "Shirt", TotalStock: 10},
		},
	}
	svc := product.NewService(mock)

	result, err := svc.GetProductList()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Shirt", result[0].Name)
}

func TestGetProductList_Error(t *testing.T) {
	mock := &mockRepo{err: assert.AnError}
	svc := product.NewService(mock)

	_, err := svc.GetProductList()
	assert.Error(t, err)
}
