package entity

import (
	"errors"
	"time"
)

var errInvalidUser = errors.New("id, email, and name are required")

type User struct {
	id        string
	email     string
	password  string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func ReconstructUser(id, email, password, name string, createdAt, updatedAt time.Time) (*User, error) {
	if id == "" || email == "" || name == "" {
		return nil, errInvalidUser
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

func (u *User) SetName(name string) {
	u.name = name
	u.updatedAt = time.Now()
}

func (u *User) SetPassword(hashedPassword string) {
	u.password = hashedPassword
	u.updatedAt = time.Now()
}
