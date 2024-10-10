package mq

import "github.com/google/uuid"

type OrderPlacedPayload struct {
	Products []*OrderPlacedProduct `json:"product"`
}

type OrderPlacedProduct struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type ProductCreatedPayload struct {
	ProductId uuid.UUID `json:"product_id"`
}

type InventoryUpdatePayload struct {
	ProductId uuid.UUID `json:"product_id"`
	Type      int32     `json:"type"`
	Number    int       `json:"number"`
}
