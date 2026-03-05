package user

type CreateUserRequest struct {
	PositionID string `json:"positionID"`
	EmployeeID string `json:"employeeID"`
	Email      string `validate:"required" json:"email"`
	Password   string `validate:"required,min=1,max=200" json:"password"`
	Role       int    `validate:"required" json:"role"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Access     bool   `json:"access"`
	Phone      string `validate:"required" json:"phone"`
}
