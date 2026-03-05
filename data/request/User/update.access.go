package user

type UpdateAccessRequest struct {
	ID     string `validate:"required,min=1,max=200" json:"id"`
	Access bool   `json:"access"`
}
