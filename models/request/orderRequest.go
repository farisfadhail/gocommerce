package request

type OrderRequest struct {
	UserOrderId uint `json:"user_order_id" validate:"required,number"`
	ProductId   uint `json:"product_id" validate:"required,number"`
	Quantity    uint `json:"quantity" validate:"required,number"`
}

type OrderUpdateRequest struct {
	Status string `json:"status" validate:"oneof=paid shipping delivered cancel"` // paid, shipping, delivered, cancel
}
