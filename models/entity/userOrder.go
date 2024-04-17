package entity

import "time"

type UserOrder struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserId     uint      `json:"user_id" gorm:"column:user_id"`
	Phone      string    `json:"phone" gorm:"column:phone"`
	Address    string    `json:"address" gorm:"column:address"`
	District   string    `json:"district" gorm:"column:district"`
	City       string    `json:"city" gorm:"column:city"`
	Province   string    `json:"province" gorm:"column:province"`
	PostalCode int       `json:"postal_code" gorm:"column:postal_code"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Order      Order     `json:"-" gorm:"-"`
}
