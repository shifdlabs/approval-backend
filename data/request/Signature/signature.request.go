package signature

type CreateSignatureRequest struct {
	UserID   string `validate:"required" json:"userId"`
	ImageURL string `validate:"required" json:"imageUrl"`
}

type UpdateSignatureRequest struct {
	ImageURL string `validate:"required" json:"imageUrl"`
}
