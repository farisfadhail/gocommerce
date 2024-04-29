package request

import "github.com/google/uuid"

type CartRequest struct {
	UserId    uuid.UUID `json:"user_id" validate:"required"`
	ProductId uint      `json:"product_id" validate:"required,number"`
	Quantity  uint      `json:"quantity" validate:"required,number"`
}

type CartUpdateRequest struct {
	Quantity uint `json:"quantity" validate:"number"`
}

type CartDeleteRequest struct {
	CartId []uint `json:"cart_id" validate:"required"`
}
