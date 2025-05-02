package user

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	Admin    Role = "admin"
	Driver   Role = "driver"
	Customer Role = "customer"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	FullName     string    `db:"full_name" json:"full_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password"`
	Role         Role      `db:"role" json:"role"`
	Phone        string    `db:"phone" json:"phone"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
