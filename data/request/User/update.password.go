package user

type UpdatePasswordRequest struct {
	ID              string `validate:"required,min=1,max=200" json:"id"`
	CurrentPassword string `validate:"min=1,max=200" json:"currentPassword"`
	NewPassword     string `validate:"required,min=1,max=200" json:"newPassword"`
}
