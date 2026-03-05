package user

type UpdateUserRequest struct {
	ID         string `json:"id"`
	PositionID string `json:"position"`
	EmployeeID string `json:"employeeID"`
	Email      string `validate:"required" json:"email"`
	Role       int    `validate:"required" json:"role"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Access     bool   `json:"access"`
	Phone      string `validate:"required" json:"phone"`
}

type UpdateUserTypeRequest struct {
	ID   string `validate:"required,min=1,max=200" json:"id"`
	Type int    `validate:"required" json:"type"`
}
