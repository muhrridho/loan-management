package entity

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateUserPayload struct {
	Email string `db:"email"`
	Name  string `db:"name"`
}

// func NewUser(email string, name string) *User {
// 	return &User{
// 		Email:     email,
// 		Name:      name,
// 		CreatedAt: time.Now(),
// 	}
// }
