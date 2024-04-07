package request

type ImageCreateRequest struct {
	ProductId uint `json:"product_id" form:"product_id" validate:"required,number"`
}
