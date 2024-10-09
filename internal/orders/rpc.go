package orders

import (
	"context"
	"log"

	"github.com/ary82/microservices/internal/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterOrdersServiceServer(grpcServer, NewUsersServer())
	return grpcServer
}

// grpc server implementation
type ordersService struct {
	proto.UnimplementedOrdersServiceServer
}

func (s *ordersService) GetOrder(context.Context, *proto.OrderId) (*proto.Order, error) {
	return &proto.Order{}, nil
}

func (s *ordersService) GetOrders(context.Context, *proto.GetOrdersParams) (*proto.OrderList, error) {
	log.Println("GetProducts")
	return &proto.OrderList{}, nil
}

func (s *ordersService) PlaceOrder(context.Context, *proto.PlaceOrderParams) (*proto.PlaceOrderResponse, error) {
	return &proto.PlaceOrderResponse{}, nil
}

func NewUsersServer() proto.OrdersServiceServer {
	s := new(ordersService)
	return s
}
