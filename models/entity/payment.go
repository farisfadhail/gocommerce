package entity

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	TransactionID   string    `json:"transaction_id" gorm:"column:transaction_id"`     // transaction_id
	Amount          int64     `json:"amount" gorm:"column:amount"`                     // grass_amount
	PaymentType     string    `json:"payment_type" gorm:"column:payment_type"`         // payment_type
	Status          string    `json:"status" gorm:"column:status"`                     // pending, success, failed (transaction_status & fraud_status)
	TransactionTime string    `json:"transaction_time" gorm:"column:transaction_time"` // transaction_time
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
