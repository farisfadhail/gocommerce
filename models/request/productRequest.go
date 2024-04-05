package request

type ProductRequest struct {
	CategoryId  int    `json:"category_id" form:"category_id" validate:"required,number"`
	Name        string `json:"name" validate:"required"`
	Price       int64  `json:"price" validate:"required,number"`
	Description string `json:"description" validate:"required"`
}
