package users

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type UsersService interface {
	GetUser(id string) (*User, error)
	GetAllUsers() ([]*User, error)
	LoginUser(LoginRequest) (*string, error)
	RegisterUser(req RegisterUserRequest) (*uuid.UUID, error)
}

type usersService struct {
	repo     UsersRepository
	producer EventProducer
}

func NewUsersService(repo UsersRepository, p EventProducer) UsersService {
	return &usersService{
		repo:     repo,
		producer: p,
	}
}

func (s *usersService) GetUser(id string) (*User, error) {
	decodeID, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.FromBytes(decodeID)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FetchUser(uuid)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:       uuid,
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

	saltDecode, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return nil, err
	}

	hashDecode, err := base64.StdEncoding.DecodeString(user.PassHash)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(req.Password), saltDecode, 1, 64*1024, 4, 32)

	match := subtle.ConstantTimeCompare(hashDecode, hash)
	if match == 0 {
		return nil, fmt.Errorf("password incorrect")
	}

	token, err := GenerateJWT(user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *usersService) RegisterUser(req RegisterUserRequest) (*uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(req.Password), salt, 1, 64*1024, 4, 32)

	base64Salt := base64.StdEncoding.EncodeToString(salt)
	base64Hash := base64.StdEncoding.EncodeToString(hash)

	err = s.repo.RegisterUser(
		id,
		req.Username,
		req.Email,
		base64Salt,
		base64Hash,
		req.UserType,
	)
	if err != nil {
		return nil, err
	}

	s.producer.Produce("USER_REGISTERED", id[:])

	return &id, nil
}
