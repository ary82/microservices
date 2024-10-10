package users

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UsersRepository interface {
	FetchUser(id uuid.UUID) (*DetailedUser, error)
	FetchUserByEmail(email string) (*DetailedUser, error)
	FetchAllUsers() ([]*User, error)
	RegisterUser(id uuid.UUID, username string, email string, salt string, passHash string, userType int32) error
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepository{
		db: db,
	}
}

const (
	fetchUserQuery        = "SELECT username, email, salt, pass_hash, user_type FROM users WHERE id = $1"
	FetchUserByEmailQuery = "SELECT id, username, salt, pass_hash, user_type FROM users WHERE email = $1"
	fetchAllUsersQuery    = "SELECT id, username, email, user_type FROM users"
	registerUserQuery     = "INSERT INTO users(id, username, email, salt, pass_hash, user_type, created_at) VALUES($1, $2, $3, $4, $5, $6, $7)"
)

func (r *usersRepository) FetchUser(id uuid.UUID) (*DetailedUser, error) {
	user := &DetailedUser{Id: id}
	err := r.db.QueryRow(fetchUserQuery, id).Scan(
		&user.Username,
		&user.Email,
		&user.Salt,
		&user.PassHash,
		&user.UserType,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FetchUserByEmail(email string) (*DetailedUser, error) {
	user := &DetailedUser{Email: email}
	err := r.db.QueryRow(FetchUserByEmailQuery, email).Scan(
		&user.Id,
		&user.Username,
		&user.Salt,
		&user.PassHash,
		&user.UserType,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FetchAllUsers() ([]*User, error) {
	users := []*User{}
	rows, err := r.db.Query(fetchAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(User)
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.UserType,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *usersRepository) RegisterUser(id uuid.UUID, username string, email string, salt string, passHash string, userType int32) error {
	_, err := r.db.Exec(
		registerUserQuery,
		id,
		username,
		email,
		salt,
		passHash,
		userType,
		time.Now().UTC(),
	)
	return err
}
