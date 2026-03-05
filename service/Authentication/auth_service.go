package authentication

import (
	model "Microservice/data/model/Authentication"
	authentication "Microservice/data/request/Authentication"
	"Microservice/helper"
)

type AuthService interface {
	Login(payload authentication.LogInRequest) (model.LoginResult, *helper.ErrorModel)
	// CheckRegisteredEmail(payload authentication.VerifyForgetPassword) bool
	// ResetPassword(payload authentication.ResetPassword) *helper.CustomError
}
