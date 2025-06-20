package warehouse_test

import (
	"ecommerce-stock-api/internal/warehouse"
	"ecommerce-stock-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransferRepo struct {
	calledWith *warehouse.TransferInput
	returnErr  error
}

func (m *mockTransferRepo) Create(w *models.Warehouse) error        { return nil }
func (m *mockTransferRepo) UpdateStatus(id uint, active bool) error { return nil }
func (m *mockTransferRepo) List() ([]models.Warehouse, error)       { return nil, nil }
func (m *mockTransferRepo) TransferStock(input warehouse.TransferInput) error {
	m.calledWith = &input
	return m.returnErr
}

func TestService_TransferStock_CallsRepo(t *testing.T) {
	mock := &mockTransferRepo{}
	svc := warehouse.NewService(mock)

	err := svc.TransferStock(1, 2, 3, 5)
	assert.NoError(t, err)
	assert.NotNil(t, mock.calledWith)
	assert.Equal(t, uint(1), mock.calledWith.ProductID)
	assert.Equal(t, uint(2), mock.calledWith.FromWarehouseID)
	assert.Equal(t, uint(3), mock.calledWith.ToWarehouseID)
	assert.Equal(t, 5, mock.calledWith.Quantity)
}
