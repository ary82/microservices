package products

import (
	"context"
	"log"

	"github.com/ary82/microservices/internal/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterProductsServiceServer(grpcServer, NewProductsServer())
	return grpcServer
}

// grpc server implementation
type productsService struct {
	proto.UnimplementedProductsServiceServer
}

func (s *productsService) GetProduct(context.Context, *proto.ProductId) (*proto.Product, error) {
	return &proto.Product{}, nil
}

func (s *productsService) GetProducts(context.Context, *proto.GetProductsParams) (*proto.ProductList, error) {
	log.Println("GetProducts")
	return &proto.ProductList{}, nil
}

func (s *productsService) AddProduct(context.Context, *proto.Product) (*proto.AddProductResponse, error) {
	return &proto.AddProductResponse{
		Uuid: []byte{},
	}, nil
}

func NewProductsServer() proto.ProductsServiceServer {
	s := new(productsService)
	return s
}
