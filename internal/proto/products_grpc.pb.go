// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.28.2
// source: products.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ProductsService_GetProduct_FullMethodName  = "/proto.ProductsService/GetProduct"
	ProductsService_GetProducts_FullMethodName = "/proto.ProductsService/GetProducts"
	ProductsService_AddProduct_FullMethodName  = "/proto.ProductsService/AddProduct"
)

// ProductsServiceClient is the client API for ProductsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductsServiceClient interface {
	GetProduct(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*Product, error)
	GetProducts(ctx context.Context, in *GetProductsParams, opts ...grpc.CallOption) (*ProductList, error)
	AddProduct(ctx context.Context, in *Product, opts ...grpc.CallOption) (*AddProductResponse, error)
}

type productsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductsServiceClient(cc grpc.ClientConnInterface) ProductsServiceClient {
	return &productsServiceClient{cc}
}

func (c *productsServiceClient) GetProduct(ctx context.Context, in *UUID, opts ...grpc.CallOption) (*Product, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Product)
	err := c.cc.Invoke(ctx, ProductsService_GetProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) GetProducts(ctx context.Context, in *GetProductsParams, opts ...grpc.CallOption) (*ProductList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProductList)
	err := c.cc.Invoke(ctx, ProductsService_GetProducts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) AddProduct(ctx context.Context, in *Product, opts ...grpc.CallOption) (*AddProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddProductResponse)
	err := c.cc.Invoke(ctx, ProductsService_AddProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductsServiceServer is the server API for ProductsService service.
// All implementations must embed UnimplementedProductsServiceServer
// for forward compatibility
type ProductsServiceServer interface {
	GetProduct(context.Context, *UUID) (*Product, error)
	GetProducts(context.Context, *GetProductsParams) (*ProductList, error)
	AddProduct(context.Context, *Product) (*AddProductResponse, error)
	mustEmbedUnimplementedProductsServiceServer()
}

// UnimplementedProductsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProductsServiceServer struct {
}

func (UnimplementedProductsServiceServer) GetProduct(context.Context, *UUID) (*Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProduct not implemented")
}
func (UnimplementedProductsServiceServer) GetProducts(context.Context, *GetProductsParams) (*ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProducts not implemented")
}
func (UnimplementedProductsServiceServer) AddProduct(context.Context, *Product) (*AddProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProduct not implemented")
}
func (UnimplementedProductsServiceServer) mustEmbedUnimplementedProductsServiceServer() {}

// UnsafeProductsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductsServiceServer will
// result in compilation errors.
type UnsafeProductsServiceServer interface {
	mustEmbedUnimplementedProductsServiceServer()
}

func RegisterProductsServiceServer(s grpc.ServiceRegistrar, srv ProductsServiceServer) {
	s.RegisterService(&ProductsService_ServiceDesc, srv)
}

func _ProductsService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UUID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductsService_GetProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).GetProduct(ctx, req.(*UUID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_GetProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductsParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).GetProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductsService_GetProducts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).GetProducts(ctx, req.(*GetProductsParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_AddProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Product)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).AddProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductsService_AddProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).AddProduct(ctx, req.(*Product))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductsService_ServiceDesc is the grpc.ServiceDesc for ProductsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ProductsService",
	HandlerType: (*ProductsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProduct",
			Handler:    _ProductsService_GetProduct_Handler,
		},
		{
			MethodName: "GetProducts",
			Handler:    _ProductsService_GetProducts_Handler,
		},
		{
			MethodName: "AddProduct",
			Handler:    _ProductsService_AddProduct_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "products.proto",
}
