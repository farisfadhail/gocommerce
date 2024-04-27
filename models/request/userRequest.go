package request

type UserUpdateRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type UserUpdateEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type UserUpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin consumer"`
}
