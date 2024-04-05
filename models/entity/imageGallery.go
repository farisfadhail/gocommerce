package entity

import "time"

type ImageGallery struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ProductId int       `json:"product_id" gorm:"column:product_id"`
	FileName  string    `json:"file_name" gorm:"column:file_name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
