package entity

import (
	"errors"
	"time"
)

var errInvalidUser = errors.New("email, hashed password, and name are required")

type User struct {
	id        uint
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
		email:     email,
		password:  hashedPassword,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructUser(id uint, email, password, name string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() uint             { return u.id }
func (u *User) Email() string         { return u.email }
func (u *User) HashedPassword() string { return u.password }
func (u *User) Name() string          { return u.name }
func (u *User) CreatedAt() time.Time  { return u.createdAt }
func (u *User) UpdatedAt() time.Time  { return u.updatedAt }
func (u *User) SetID(id uint)         { u.id = id }
