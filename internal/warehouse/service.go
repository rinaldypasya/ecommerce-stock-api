package warehouse

import "ecommerce-stock-api/models"

type Service interface {
	Create(name string, shopID uint) error
	UpdateStatus(id uint, active bool) error
	List() ([]models.Warehouse, error)
	TransferStock(productID, fromWarehouseID, toWarehouseID uint, quantity int) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(name string, shopID uint) error {
	return s.repo.Create(&models.Warehouse{Name: name, ShopID: shopID, IsActive: true})
}

func (s *service) UpdateStatus(id uint, active bool) error {
	return s.repo.UpdateStatus(id, active)
}

func (s *service) List() ([]models.Warehouse, error) {
	return s.repo.List()
}

func (s *service) TransferStock(productID, fromWarehouseID, toWarehouseID uint, quantity int) error {
	input := TransferInput{
		ProductID:       productID,
		FromWarehouseID: fromWarehouseID,
		ToWarehouseID:   toWarehouseID,
		Quantity:        quantity,
	}
	return s.repo.TransferStock(input)
}
