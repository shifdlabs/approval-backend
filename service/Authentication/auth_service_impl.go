package authentication

import (
	model "Microservice/data/model/Authentication"
	authentication "Microservice/data/request/Authentication"
	"Microservice/helper"
	userRepository "Microservice/repository/User"

	"github.com/go-playground/validator/v10"
)

type AuthServiceImpl struct {
	UserRepository userRepository.UserRepository
	Validate       *validator.Validate
}

func NewAuthServiceImpl(userRepository userRepository.UserRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (t AuthServiceImpl) Login(payload authentication.LogInRequest) (model.LoginResult, *helper.ErrorModel) {
	user, err := t.UserRepository.GetByEmail(payload.Email)
	if err != nil {
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, err
	}

	errVerifyPassword := helper.VerifyPassword(user.Password, payload.Password)
	if errVerifyPassword != nil {
		msg := "Incorrect password"
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, helper.ErrorCatcher(errVerifyPassword, 400, &msg)
	}

	accessToken, _ := helper.GenerateAccessToken(user.ID.String())
	refreshToken, _ := helper.GenerateRefreshToken(user.ID.String())

	return model.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// func (t AuthServiceImpl) CheckRegisteredEmail(payload authentication.VerifyForgetPassword) bool {
// 	user, _ := t.UserRepository.GetByEmail(payload.Email)

// 	if user != nil {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (t AuthServiceImpl) ResetPassword(payload authentication.ResetPassword) *helper.CustomError {
// 	err := t.UserRepository.UpdatePasssword(payload.Email, payload.NewPassword)

// 	if err != nil {
// 		return err
// 	} else {
// 		return nil
// 	}
// }
