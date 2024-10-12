package products

import (
	"context"

	"github.com/ary82/microservices/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string, ps ProductsService) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterProductsServiceServer(grpcServer, NewProductsServer(ps))
	return grpcServer
}

// grpc server implementation
type productDomainRpc struct {
	service ProductsService

	// implement grpc
	proto.UnimplementedProductsServiceServer
}

func NewProductsServer(ps ProductsService) proto.ProductsServiceServer {
	return &productDomainRpc{
		service: ps,
	}
}

func (s *productDomainRpc) GetProduct(ctx context.Context, in *proto.ProductId) (*proto.Product, error) {
	uuid, err := uuid.FromBytes(in.Value)
	if err != nil {
		return nil, err
	}

	product, err := s.service.GetProduct(uuid)
	if err != nil {
		return nil, err
	}
	return &proto.Product{
		Name:        product.Name,
		Description: product.Desc,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *productDomainRpc) GetProducts(context.Context, *proto.GetProductsParams) (*proto.ProductList, error) {
	products, err := s.service.GetProducts()
	if err != nil {
		return nil, err
	}

	protoProducts := new(proto.ProductList)
	protoProducts.Number = int64(len(products))

	for _, v := range products {
		protoProducts.Products = append(protoProducts.Products, &proto.Product{
			Id:          v.Id[:],
			Name:        v.Name,
			Description: v.Desc,
			Price:       v.Price,
			Stock:       v.Stock,
		})
	}

	return protoProducts, nil
}

func (s *productDomainRpc) AddProduct(ctx context.Context, in *proto.AddProductRequest) (*proto.AddProductResponse, error) {
	id, err := s.service.AddProduct(AddProductRequest{
		Name:  in.Name,
		Desc:  in.Description,
		Price: in.Price,
		Stock: in.Stock,
	})
	if err != nil {
		return nil, err
	}

	return &proto.AddProductResponse{
		Uuid: (*id)[:],
	}, nil
}

func (s *productDomainRpc) UpdateInventory(ctx context.Context, in *proto.UpdateInventoryRequest) (*proto.UpdateInventoryResponse, error) {
	uuid, err := uuid.FromBytes(in.Id)
	if err != nil {
		return nil, err
	}
	err = s.service.UpdateInventory(UpdateInventoryRequest{
		ProductId: uuid,
		Number:    in.Number,
		Type:      int32(in.Type),
	})
	if err != nil {
		return nil, err
	}

	return &proto.UpdateInventoryResponse{}, nil
}
