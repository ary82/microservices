package users

import (
	"context"

	"github.com/ary82/microservices/internal/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(port string, s UsersService) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, NewUsersServer(s))
	return grpcServer
}

// grpc server implementation
type usersServiceRpc struct {
	service UsersService

	// implement grpc
	proto.UnimplementedUserServiceServer
}

func NewUsersServer(s UsersService) proto.UserServiceServer {
	return &usersServiceRpc{
		service: s,
	}
}

func (s *usersServiceRpc) GetUser(ctx context.Context, in *proto.UserId) (*proto.User, error) {
	user, err := s.service.GetUser(in.Value)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:       user.Id[:],
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
	protoUsers.Number = int64(len(users))
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

func (s *usersServiceRpc) RegisterUser(ctx context.Context, in *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	req := RegisterUserRequest{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
		UserType: int32(in.Type),
	}
	id, err := s.service.RegisterUser(req)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterUserResponse{
		Id: id[:],
	}, nil
}
