package entity

import "time"

type Category struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	Slug      string    `json:"slug" gorm:"column:slug;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Products  []Product `json:"-" gorm:"-"`
}
