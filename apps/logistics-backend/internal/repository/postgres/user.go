package postgres

import (
	"logistics-backend/internal/domain/user"

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
		INSERT INTO users (full_name, email, password_hash, role, phone, created_at)
		VALUES (:full_name, :email, :password_hash, :role, :phone, NOW())
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&u.ID, u)
}

func (r *UserRepository) GetByID(id int64) (*user.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var u user.User
	err := r.db.Get(&u, query, id)
	return &u, err
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	var u user.User
	err := r.db.Get(&u, query, email)
	return &u, err
}

func (r *UserRepository) List() ([]*user.User, error) {
	query := `SELECT * FROM users ORDER BY created_at DESC`
	var users []*user.User
	err := r.db.Select(&users, query)
	return users, err
}
