package request

type UserRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Username string `json:"username" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type UserUpdateRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type UserUpdateEmailRequest struct {
	Email string `json:"email" validate:"required"`
}
