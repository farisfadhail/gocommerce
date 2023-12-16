package entity

import "time"

type User struct {
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"size:100;column:username"`
	Email     string    `gorm:"size:255;column:email"`
	Password  string    `gorm:"column:password"`
	Role      string    `gorm:"column:role"` // user & admin
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
