package users

import (
	"context"
	"log"

	"github.com/ary82/microservices/internal/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, NewUsersServer())
	return grpcServer
}

// grpc server implementation
type usersService struct {
	proto.UnimplementedUserServiceServer
}

func (s *usersService) GetUser(context.Context, *proto.UserId) (*proto.User, error) {
	return &proto.User{}, nil
}

func (s *usersService) GetUsers(context.Context, *proto.GetUsersParams) (*proto.Users, error) {
	log.Println("GetProducts")
	return &proto.Users{}, nil
}

func (s *usersService) Login(context.Context, *proto.LoginRequest) (*proto.LoginResponse, error) {
	return &proto.LoginResponse{}, nil
}

func NewUsersServer() proto.UserServiceServer {
	s := new(usersService)
	return s
}
