package request

type PaymentRequest struct {
	PaymentNumber string  `json:"payment_number" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,number"`
	Status        string  `json:"status" validate:"required"`
	PaymentType   string  `json:"payment_type" validate:"required"`
}

type PaymentUpdateStatusRequest struct {
	Status string `json:"status" validate:"oneof=Success Failed"`
}
