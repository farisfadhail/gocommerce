package request

type UserOrderRequest struct {
	UserId     uint   `json:"user_id" validate:"required,number"`
	Address    string `json:"address" validate:"required"`
	District   string `json:"district" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode int    `json:"postal_code" validate:"required,number"`
}

type UserOrderUpdateRequest struct {
	UserId     uint   `json:"user_id" validate:"number"`
	Address    string `json:"address"`
	District   string `json:"district"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode int    `json:"postal_code" validate:"number"`
}
