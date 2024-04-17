package request

type OrderRequest struct {
	// order request
	ProductId uint `json:"product_id" validate:"required,number"`
	Quantity  uint `json:"quantity" validate:"required,number"`
	// user order request
	UserId     uint   `json:"user_id" validate:"required,number"`
	Phone      string `json:"phone" validate:"required"`
	Address    string `json:"address" validate:"required"`
	District   string `json:"district" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode int    `json:"postal_code" validate:"required,number"`
	// payment request
	PaymentType string `json:"payment_type" validate:"required"`
	TokenID     string `json:"token_id"`
}

type OrderByCartRequest struct {
	// order request
	CartId []uint `json:"cart_id" validate:"required"`
	// user order request
	UserId     uint   `json:"user_id" validate:"required,number"`
	Phone      string `json:"phone" validate:"required,number"`
	Address    string `json:"address" validate:"required"`
	District   string `json:"district" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode int    `json:"postal_code" validate:"required,number"`
	// payment request
	PaymentType string `json:"payment_type" validate:"required"`
	TokenID     string `json:"token_id"`
}

type OrderUpdateRequest struct {
	Status string `json:"status" validate:"oneof=Paid Shipping Delivered Canceled"` // paid, shipping, delivered, canceled
}
