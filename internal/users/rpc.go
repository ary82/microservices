package users

import (
	"context"

	"github.com/ary82/microservices/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, NewUsersServer())
	return grpcServer
}

// grpc server implementation
type usersServiceRpc struct {
	service UsersService

	// implement grpc
	proto.UnimplementedUserServiceServer
}

func (s *usersServiceRpc) GetUser(ctx context.Context, in *proto.UserId) (*proto.User, error) {
	id, err := uuid.FromBytes(in.Value)
	if err != nil {
		return nil, err
	}
	user, err := s.service.GetUser(id)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:       id[:],
		Email:    user.Email,
		Username: user.Username,
		Type:     proto.UserType(user.UserType),
	}, nil
}

func (s *usersServiceRpc) GetUsers(context.Context, *proto.GetUsersParams) (*proto.Users, error) {
	users, err := s.service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	protoUsers := new(proto.Users)
	for _, v := range users {
		protoUsers.Users = append(protoUsers.Users, &proto.User{
			Id:       v.Id[:],
			Email:    v.Email,
			Username: v.Username,
			Type:     proto.UserType(v.UserType),
		})
	}

	return protoUsers, nil
}

func (s *usersServiceRpc) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.service.LoginUser(LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{
		AccessToken: *token,
	}, nil
}

func NewUsersServer() proto.UserServiceServer {
	s := new(usersServiceRpc)
	return s
}
