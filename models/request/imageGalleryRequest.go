package request

type ImageCreateRequest struct {
	ProductId int `json:"product_id" form:"product_id" validate:"required,number"`
}
