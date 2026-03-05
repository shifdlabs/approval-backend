package user

type UpdateEmailRequest struct {
	NewEmail string `validate:"required" json:"newEmail"`
}
