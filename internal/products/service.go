package products

import (
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductsService interface {
	GetProduct(id uuid.UUID) (*DetailedProduct, error)
	GetProducts() ([]*DetailedProduct, error)
	AddProduct(req AddProductRequest) (*uuid.UUID, error)
	UpdateInventory(req UpdateInventoryRequest) error
}

type productsService struct {
	repo   ProductsRepository
	mqChan *amqp.Channel
}

func NewProductsService(repo ProductsRepository, ch *amqp.Channel) ProductsService {
	return &productsService{
		repo:   repo,
		mqChan: ch,
	}
}

func (p *productsService) GetProduct(id uuid.UUID) (*DetailedProduct, error) {
	product, err := p.repo.FetchProduct(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productsService) GetProducts() ([]*DetailedProduct, error) {
	products, err := p.repo.FetchAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productsService) AddProduct(req AddProductRequest) (*uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	err = p.repo.AddProduct(
		id,
		req.Name,
		req.Desc,
		req.Price,
		req.Stock,
	)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (p *productsService) UpdateInventory(req UpdateInventoryRequest) error {
	switch req.Type {
	case 1:
		err := p.repo.UpdateInventoryAdd(req.ProductId, req.Number)
		return err
	case 2:
		err := p.repo.UpdateInventorySubtract(req.ProductId, req.Number)
		return err
	case 3:
		err := p.repo.UpdateInventoryDelete(req.ProductId)
		return err
	default:
		return fmt.Errorf("update inventory type is invalid")
	}
}
