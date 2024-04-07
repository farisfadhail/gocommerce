package entity

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FullName  string    `json:"full_name" gorm:"size:255;column:full_name"`
	Username  string    `json:"username" gorm:"size:100;column:username;uniqueIndex"`
	Email     string    `json:"email" gorm:"size:255;column:email"`
	Password  string    `json:"-" gorm:"column:password"`
	Role      string    `json:"role" gorm:"column:role"` // consumer & admin
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
