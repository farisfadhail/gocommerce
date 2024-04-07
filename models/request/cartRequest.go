package request

type CartRequest struct {
	UserId    uint `json:"user_id" validate:"required,number"`
	ProductId uint `json:"product_id" validate:"required,number"`
	Quantity  uint `json:"quantity" validate:"required,number"`
}

type CartUpdateRequest struct {
	Quantity uint `json:"quantity" validate:"number"`
}

type CartDeleteRequest struct {
	CartId []uint `json:"cart_id" validate:"required"`
}
