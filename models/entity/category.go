package entity

type Category struct {
	ID       int       `gorm:"column:id;primaryKey"`
	Name     string    `gorm:"column:name"`
	Slug     string    `gorm:"column:slug"`
	Products []Product `gorm:"foreignKey:category_id;references:id"`
}
