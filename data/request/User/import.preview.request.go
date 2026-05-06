package user

type PreviewImportRequest struct {
	ColumnMapping map[string]string `json:"columnMapping"` // maps field name to Excel column name
	// Example: {"email": "Email", "firstName": "First Name", "lastName": "Last Name", ...}
}

type ImportedUserData struct {
	EmployeeID string `json:"employeeID"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Role       int    `json:"role"`
	Phone      string `json:"phone"`
	PositionID string `json:"positionID"`
	Password   string `json:"password"` // Optional, will use default if empty
}
