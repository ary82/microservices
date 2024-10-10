package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ary82/microservices/internal/mq"
	"github.com/google/uuid"
)

type ProductsService interface {
	GetProduct(id uuid.UUID) (*DetailedProduct, error)
	GetProducts() ([]*DetailedProduct, error)
	AddProduct(req AddProductRequest) (*uuid.UUID, error)
	UpdateInventory(req UpdateInventoryRequest) error
}

type productsService struct {
	repo     ProductsRepository
	producer EventProducer
}

func NewProductsService(repo ProductsRepository, p EventProducer) ProductsService {
	return &productsService{
		repo:     repo,
		producer: p,
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

	// Send Message
	payload := mq.ProductCreatedPayload{
		ProductId: id,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(payload)
	if err != nil {
		log.Println(err)
	}
	p.producer.Produce("PRODUCT_CREATED", b.Bytes())

	return &id, nil
}

func (p *productsService) UpdateInventory(req UpdateInventoryRequest) error {
	switch req.Type {
	case 1:
		err := p.repo.UpdateInventoryAdd(req.ProductId, req.Number)
		if err != nil {
			return err
		}
	case 2:
		err := p.repo.UpdateInventorySubtract(req.ProductId, req.Number)
		if err != nil {
			return err
		}
	case 3:
		err := p.repo.UpdateInventoryDelete(req.ProductId)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("update inventory type is invalid")
	}

	// Send Message
	payload := mq.InventoryUpdatePayload{
		ProductId: req.ProductId,
		Type:      req.Type,
		Number:    int(req.Number),
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(payload)
	if err != nil {
		log.Println(err)
	}
	p.producer.Produce("INVENTORY_UPDATE", b.Bytes())

	return nil
}
