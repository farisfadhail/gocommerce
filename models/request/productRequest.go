package request

type ProductRequest struct {
	CategoryId  int    `json:"category_id" form:"category_id" validate:"required,number"`
	Name        string `json:"name" validate:"required"`
	Price       uint64 `json:"price" validate:"required,number"`
	Description string `json:"description" validate:"required"`
	Quantity    uint64 `json:"quantity" validate:"required,number"`
}

type ProductUpdateRequest struct {
	CategoryId  int    `json:"category_id" form:"category_id" validate:"number"`
	Name        string `json:"name"`
	Price       uint64 `json:"price" validate:"number"`
	Description string `json:"description"`
	Quantity    uint64 `json:"quantity" validate:"number"`
}
