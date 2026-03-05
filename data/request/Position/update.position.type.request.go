package position

type UpdatePositionRequest struct {
	ID   string `validate:"required,min=1,max=200" json:"id"`
	Name string `validate:"required,min=1,max=200" json:"name"`
}
