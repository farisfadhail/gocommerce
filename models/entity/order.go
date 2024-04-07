package entity

import "time"

type Order struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserOrderId uint      `json:"user_order_id" gorm:"column:user_order_id"`
	ProductId   uint      `json:"product_id" gorm:"column:product_id"`
	Quantity    uint      `json:"quantity" gorm:"column:quantity"`
	Status      string    `json:"status" gorm:"column:status"` // pending, paid, shipping, delivered, cancel
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
