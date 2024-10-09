package users

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Username string
	Email    string
	UserType int32
}

type DetailedUser struct {
	Id       uuid.UUID
	Username string
	Email    string
	Salt     string
	PassHash string
	UserType int32
}

type LoginRequest struct {
	Email    string
	Password string
}

type RegisterUserRequest struct {
	Username string
	Email    string
	Password string
	UserType int32
}
