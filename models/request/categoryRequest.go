package request

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name"`
}
