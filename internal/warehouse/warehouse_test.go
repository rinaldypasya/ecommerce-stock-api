package warehouse_test

import (
	"ecommerce-stock-api/internal/warehouse"
	"ecommerce-stock-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	failCreate  bool
	failUpdate  bool
	failList    bool
	statusCalls []bool
	warehouses  []models.Warehouse
}

// TransferStock implements warehouse.Repository.
func (m *mockRepo) TransferStock(input warehouse.TransferInput) error {
	panic("unimplemented")
}

func (m *mockRepo) Create(w *models.Warehouse) error {
	if m.failCreate {
		return assert.AnError
	}
	m.warehouses = append(m.warehouses, *w)
	return nil
}

func (m *mockRepo) UpdateStatus(id uint, active bool) error {
	if m.failUpdate {
		return assert.AnError
	}
	m.statusCalls = append(m.statusCalls, active)
	return nil
}

func (m *mockRepo) List() ([]models.Warehouse, error) {
	if m.failList {
		return nil, assert.AnError
	}
	return m.warehouses, nil
}

func TestCreateWarehouse_Success(t *testing.T) {
	mock := &mockRepo{}
	svc := warehouse.NewService(mock)

	err := svc.Create("Gudang 1", 1)
	assert.NoError(t, err)
	assert.Len(t, mock.warehouses, 1)
	assert.Equal(t, "Gudang 1", mock.warehouses[0].Name)
}

func TestUpdateStatus_Success(t *testing.T) {
	mock := &mockRepo{}
	svc := warehouse.NewService(mock)

	err := svc.UpdateStatus(1, false)
	assert.NoError(t, err)
	assert.Len(t, mock.statusCalls, 1)
	assert.False(t, mock.statusCalls[0])
}
