package entity

import "time"

type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id" gorm:"column:user_id"`
	ProductId uint      `json:"product_id" gorm:"column:product_id"`
	Quantity  uint      `json:"quantity" gorm:"column:quantity"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
