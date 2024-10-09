package users

import (
	"crypto/subtle"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type UsersService interface {
	GetUser(id uuid.UUID) (*User, error)
	GetAllUsers() ([]*User, error)
	LoginUser(LoginRequest) (*string, error)
}

type usersService struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}

func (s *usersService) GetUser(id uuid.UUID) (*User, error) {
	user, err := s.repo.FetchUser(id)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:       id,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
	}, nil
}

func (s *usersService) GetAllUsers() ([]*User, error) {
	users, err := s.repo.FetchAllUsers()
	if err != nil {
		return nil, err
	}
	return users, err
}

func (s *usersService) LoginUser(req LoginRequest) (*string, error) {
	user, err := s.repo.FetchUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey([]byte(req.Password), []byte(user.Salt), 1, 64*1024, 4, 32)

	match := subtle.ConstantTimeCompare([]byte(req.Password), hash)
	if match == 0 {
		return nil, fmt.Errorf("password incorrect")
	}

	token, err := GenerateJWT(user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
