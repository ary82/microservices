package orders

import (
	"time"

	"github.com/google/uuid"
)

type DetailedOrder struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	TotalQuantity int
	TotalPrice    int
	CreatedAt     time.Time
	Products      []*DetailedOrderProduct
}

type Order struct {
	Id            uuid.UUID
	UserId        uuid.UUID
	TotalQuantity int
	TotalPrice    int
	CreatedAt     time.Time
}

type DetailedOrderProduct struct {
	ProductId uuid.UUID
	Quantity  int
	Price     int
}

type PlaceOrderRequest struct {
	UserId   uuid.UUID
	Products []*DetailedOrderProduct
}
