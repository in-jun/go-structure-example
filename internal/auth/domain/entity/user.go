package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	errInvalidUser            = errors.New("email, hashed password, and name are required")
	errInvalidReconstructUser = errors.New("id, email, hashed password, and name are required")
)

type User struct {
	id        string
	email     string
	password  string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email, hashedPassword, name string) (*User, error) {
	if email == "" || hashedPassword == "" || name == "" {
		return nil, errInvalidUser
	}
	now := time.Now()
	return &User{
		id:        uuid.New().String(),
		email:     email,
		password:  hashedPassword,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructUser(id, email, password, name string, createdAt, updatedAt time.Time) (*User, error) {
	if id == "" || email == "" || password == "" || name == "" {
		return nil, errInvalidReconstructUser
	}
	return &User{
		id:        id,
		email:     email,
		password:  password,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (u *User) ID() string             { return u.id }
func (u *User) Email() string          { return u.email }
func (u *User) HashedPassword() string { return u.password }
func (u *User) Name() string           { return u.name }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }
