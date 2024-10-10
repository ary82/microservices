package products

import (
	"time"

	"github.com/google/uuid"
)

type DetailedProduct struct {
	Id        uuid.UUID
	Name      string
	Desc      string
	Price     int64
	Stock     int64
	CreatedAt time.Time
}

type AddProductRequest struct {
	Name  string
	Desc  string
	Price int64
	Stock int64
}

type UpdateInventoryRequest struct {
	ProductId uuid.UUID
	Number    int64
	Type      int32
}
