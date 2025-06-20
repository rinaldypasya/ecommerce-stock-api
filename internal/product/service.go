package product

type Service interface {
	GetProductList() ([]ProductWithStock, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetProductList() ([]ProductWithStock, error) {
	products, err := s.repo.ListProductsWithAvailableStock()
	if err != nil {
		return nil, err
	}
	return products, nil
}
