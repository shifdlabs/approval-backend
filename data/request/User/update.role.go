package user

type UpdateRoleRequest struct {
	ID   string `validate:"required,min=1,max=200" json:"id"`
	Role int    `validate:"required" json:"role"`
}
