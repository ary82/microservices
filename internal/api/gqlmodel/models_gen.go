// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

type AllOrdersOrder struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	PriceTotal int    `json:"price_total"`
	Quantity   int    `json:"quantity"`
}

type CreateProduct struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mutation struct {
}

type NewUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}

type Order struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	PriceTotal    int             `json:"price_total"`
	Quantity      int             `json:"quantity"`
	OrderProducts []*OrderProduct `json:"order_products"`
}

type OrderProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type PlaceOrder struct {
	OrderProducts []*PlaceOrderProduct `json:"order_products"`
}

type PlaceOrderProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
}

type Query struct {
}

type Token struct {
	Token string `json:"token"`
}

type UpdateInventory struct {
	ID     string `json:"id"`
	Type   int    `json:"type"`
	Number int    `json:"number"`
}

type UpdateInventoryResponse struct {
	Success bool `json:"success"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Type     int    `json:"type"`
}
