package user

type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` //raw password from client
	Role     Role   `json:"role" binding:"required,oneof=admin driver customer"`
	Phone    string `json:"phone" binding:"required"`
}

func (r *CreateUserRequest) ToUser() *User {
	return &User{
		FullName:     r.FullName,
		Email:        r.Email,
		PasswordHash: r.Password,
		Role:         r.Role,
		Phone:        r.Phone,
	}
}
