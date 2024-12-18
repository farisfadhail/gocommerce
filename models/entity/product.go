package entity

import "time"

type Product struct {
	ID             uint64         `json:"id" gorm:"column:id;primaryKey"`
	CategoryId     uint           `json:"category_id" gorm:"column:category_id"`
	Name           string         `json:"name" gorm:"column:name"`
	Slug           string         `json:"slug" gorm:"column:slug;uniqueIndex"`
	Price          int64          `json:"price" gorm:"column:price"`
	Description    string         `json:"description" gorm:"column:description"`
	Quantity       uint64         `json:"quantity" gorm:"column:quantity"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	ImageGalleries []ImageGallery `json:"-" gorm:"-"`
	Cart           []Cart         `json:"-" gorm:"-"`
	Order          []Order        `json:"-" gorm:"-"`
}
