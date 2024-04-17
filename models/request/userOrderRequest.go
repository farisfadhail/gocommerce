package request

type UserOrderUpdateRequest struct {
	UserId     uint   `json:"user_id" validate:"number"`
	Phone      string `json:"phone" validate:"number"`
	Address    string `json:"address"`
	District   string `json:"district"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode int    `json:"postal_code" validate:"number"`
}
