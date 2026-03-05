package authentication

type LogInRequest struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"  validate:"required"`
}
