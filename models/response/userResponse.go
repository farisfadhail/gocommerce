package response

import "time"

type UserResponse struct {
	ID        uint      `gorm:"primaryKey"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-" gorm:"column:password"`
	Role      string    `json:"role"` // user & admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateEmailUserResponse struct {
	Email     string    `json:"email" gorm:"uniqueIndex"`
	UpdatedAt time.Time `json:"updatedAt"`
}
