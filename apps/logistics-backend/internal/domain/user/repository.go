package user

type Repository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	List() ([]*User, error)
}
