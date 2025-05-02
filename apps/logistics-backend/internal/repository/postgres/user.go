package postgres

import (
	"logistics-backend/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, phone)
		VALUES (:full_name, :email, :password_hash, :role, :phone)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&u.ID, u)
}

func (r *UserRepository) GetByID(id uuid.UUID) (*user.User, error) {
	query := `SELECT id, full_name, email, password_hash, role, phone FROM users WHERE id = $1`
	var u user.User
	err := r.db.Get(&u, query, id)
	return &u, err
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	query := `SELECT id, full_name, email, password_hash, role, phone FROM users WHERE email = $1`
	var u user.User
	err := r.db.Get(&u, query, email)
	return &u, err
}

func (r *UserRepository) List() ([]*user.User, error) {
	query := `SELECT id, full_name, email FROM users`
	var users []*user.User
	err := r.db.Select(&users, query)
	return users, err
}
