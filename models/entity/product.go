package entity

type Product struct {
	ID          int    `gorm:"column:id;primaryKey"`
	CategoryId  int    `gorm:"column:category_id"`
	Name        string `gorm:"column:name"`
	Slug        string `gorm:"column:slug"`
	Price       int    `gorm:"column:price"`
	Description string `gorm:"column:description"`
}
