package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	OrderNumber string    `json:"order_number" gorm:"column:order_number;uniqueIndex"`
	UserOrderId uint      `json:"user_order_id" gorm:"column:user_order_id"`
	ProductId   uint      `json:"product_id" gorm:"column:product_id"`
	Quantity    uint      `json:"quantity" gorm:"column:quantity"`
	TotalPrice  uint64    `json:"total_price" gorm:"column:total_price"`
	Status      string    `json:"status" gorm:"column:status"` // pending, paid, shipping, delivered, cancel
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
