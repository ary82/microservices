package orders

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/ary82/microservices/internal/mq"
	"github.com/google/uuid"
)

type OrdersService interface {
	GetOrder(id uuid.UUID) (*DetailedOrder, error)
	GetAllOrders() ([]*Order, error)
	PlaceOrder(req PlaceOrderRequest) (*uuid.UUID, error)
}

type ordersService struct {
	repo     OrdersRepository
	producer EventProducer
}

func NewOrdersService(repo OrdersRepository, p EventProducer) OrdersService {
	return &ordersService{
		repo:     repo,
		producer: p,
	}
}

func (s *ordersService) GetOrder(id uuid.UUID) (*DetailedOrder, error) {
	order, err := s.repo.FetchOrder(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *ordersService) GetAllOrders() ([]*Order, error) {
	orders, err := s.repo.FetchAllOrders()
	return orders, err
}

func (s *ordersService) PlaceOrder(req PlaceOrderRequest) (*uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	totalQuantity := 0
	totalPrice := 0
	for _, v := range req.Products {
		totalQuantity += v.Quantity
		totalPrice += v.Price * v.Quantity
	}

	err = s.repo.CreateOrder(
		id,
		req.UserId,
		totalQuantity,
		totalPrice,
		req.Products,
	)
	if err != nil {
		return nil, err
	}

	// Send Message
	payload := mq.OrderPlacedPayload{}
	for _, v := range req.Products {
		payload.Products = append(payload.Products, &mq.OrderPlacedProduct{
			ProductId: v.ProductId,
			Quantity:  v.Quantity,
		})
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(payload)
	if err != nil {
		log.Println(err)
	}
	s.producer.Produce("ORDER_PLACED", b.Bytes())

	return &id, nil
}
