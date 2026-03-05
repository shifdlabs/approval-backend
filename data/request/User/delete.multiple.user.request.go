package user

type DeleteMultipleUserRequest struct {
	IDs []string `validate:"required" json:"ids"`
}
