package request

import "github.com/google/uuid"

type UserOrderUpdateRequest struct {
	UserId     uuid.UUID `json:"user_id"`
	Phone      string    `json:"phone" validate:"number"`
	Address    string    `json:"address"`
	District   string    `json:"district"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode int       `json:"postal_code" validate:"number"`
}
