package user

type UpdateBiodataRequest struct {
	PositionID string `json:"position"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Phone      string `validate:"required" json:"phone"`
}
