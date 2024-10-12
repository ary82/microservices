package orders

import (
	"context"

	"github.com/ary82/microservices/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string, os OrdersService) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterOrdersServiceServer(grpcServer, NewOrdersServer(os))
	return grpcServer
}

// grpc server implementation
type ordersDomainRpc struct {
	service OrdersService

	// implement orders grpc
	proto.UnimplementedOrdersServiceServer
}

func NewOrdersServer(s OrdersService) proto.OrdersServiceServer {
	return &ordersDomainRpc{
		service: s,
	}
}

func (s *ordersDomainRpc) GetOrder(ctx context.Context, in *proto.OrderId) (*proto.Order, error) {
	id, err := uuid.FromBytes(in.Value)
	if err != nil {
		return nil, err
	}

	order, err := s.service.GetOrder(id)
	if err != nil {
		return nil, err
	}

	resp := &proto.Order{
		UserId:     order.UserId[:],
		PriceTotal: int64(order.TotalPrice),
		Quantity:   int64(order.TotalQuantity),
	}

	for _, v := range order.Products {
		resp.Products = append(resp.Products, &proto.OrderProduct{
			ProductId: v.ProductId[:],
			Quantity:  int64(v.Quantity),
			Price:     int64(v.Price),
		})
	}
	return resp, nil
}

func (s *ordersDomainRpc) GetOrders(context.Context, *proto.GetOrdersParams) (*proto.OrderList, error) {
	orders, err := s.service.GetAllOrders()
	if err != nil {
		return nil, err
	}
	resp := &proto.OrderList{
		Number: int64(len(orders)),
	}

	for _, v := range orders {
		resp.Orders = append(resp.Orders, &proto.OrderListOrder{
			OrderId:    v.Id[:],
			UserId:     v.UserId[:],
			PriceTotal: int64(v.TotalPrice),
			Quantity:   int64(v.TotalQuantity),
		})
	}
	return resp, nil
}

func (s *ordersDomainRpc) PlaceOrder(ctx context.Context, in *proto.PlaceOrderParams) (*proto.PlaceOrderResponse, error) {
	userId, err := uuid.FromBytes(in.UserId)
	if err != nil {
		return nil, err
	}
	req := PlaceOrderRequest{
		UserId: userId,
	}
	for _, v := range in.Products {
		productId, err := uuid.FromBytes(v.ProductId)
		if err != nil {
			return nil, err
		}

		req.Products = append(req.Products, &DetailedOrderProduct{
			ProductId: productId,
			Quantity:  int(v.Quantity),
			Price:     int(v.Price),
		})
	}

	orderId, err := s.service.PlaceOrder(req)
	if err != nil {
		return nil, err
	}

	return &proto.PlaceOrderResponse{
		Uuid: (*orderId)[:],
	}, nil
}
