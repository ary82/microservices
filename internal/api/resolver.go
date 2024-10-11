package api

import "github.com/ary82/microservices/internal/proto"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService    proto.UserServiceClient
	ProductService proto.ProductsServiceClient
	OrderService   proto.OrdersServiceClient
}
