package orders

import (
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrdersService interface {
	GetOrder(id uuid.UUID) (*DetailedOrder, error)
	GetAllOrders() ([]*Order, error)
	PlaceOrder(req PlaceOrderRequest) (*uuid.UUID, error)
}

type ordersService struct {
	repo   OrdersRepository
	mqChan *amqp.Channel
}

func NewOrdersService(repo OrdersRepository, ch *amqp.Channel) OrdersService {
	return &ordersService{
		repo:   repo,
		mqChan: ch,
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

	return &id, nil
}
