package controller

import (
	authentication "Microservice/data/request/Authentication"
	user "Microservice/data/request/User"
	response "Microservice/data/response"
	userResponse "Microservice/data/response/User"
	"Microservice/helper"
	"Microservice/model"
	service "Microservice/service/Authentication"
	userService "Microservice/service/User"
	"Microservice/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
	userService userService.UserService
}

func NewAuthController(service service.AuthService, userService userService.UserService) *AuthController {
	return &AuthController{authService: service, userService: userService}
}

func (controller *AuthController) LogIn(ctx *gin.Context) {
	var payload *authentication.LogInRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Invalid Structure"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	loginResult, err := controller.authService.Login(*payload)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	}

	// contextID := context.TODO()
	// now := time.Now()

	// errAccess := config.RedisClient.Set(contextID, loginResult.AccessToken.Identifier, loginResult.AccessToken.UserID, time.Unix(*loginResult.AccessToken.ExpiresIn, 0).Sub(now)).Err()
	// if errAccess != nil {
	// 	msg := "Internal Server Error"
	// 	utils.ErrorResponse(ctx, helper.ErrorModel{Code: 500, Message: msg})
	// }

	// errRefresh := config.RedisClient.Set(contextID, loginResult.RefreshToken.Identifier, loginResult.RefreshToken.UserID, time.Unix(*loginResult.RefreshToken.ExpiresIn, 0).Sub(now)).Err()
	// if errRefresh != nil {
	// 	msg := "Internal Server Error"
	// 	utils.ErrorResponse(ctx, helper.ErrorModel{Code: 500, Message: msg})
	// }

	// Set Cookie Here with Gin
	// ctx.SetCookie("access_token", *loginResult.AccessToken.Token, 3600, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", *loginResult.RefreshToken.Token, 3600, "/", "localhost", false, true)

	// ctx.SetSameSite(http.SameSiteNoneMode)
	// ctx.SetCookie("refreshToken", loginResult.RefreshToken, int(helper.RefreshTTL.Seconds()), "/", "", false, true)

	utils.SuccessResponse(ctx,
		userResponse.LoginResponse{
			AccessToken:      loginResult.AccessToken,
			RefreshToken:     loginResult.RefreshToken,
			UserAbilityRules: controller.GetUserAbilityRules(loginResult.User.Role),
			Id:               loginResult.User.ID.String(),
			Access:           loginResult.User.Access,
			Name:             loginResult.User.FirstName + " " + loginResult.User.LastName,
			Role:             loginResult.User.Role,
			JobPosition:      getPositionName(loginResult.User.Position),
		})
}

// For Testing Only, Please Remove this API for Production
func (controller *AuthController) Register(ctx *gin.Context) {
	var payload *user.CreateUserRequest

	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Structure Error"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 500, Message: msg})
	}

	errResult := controller.userService.Create(*payload)
	if errResult != nil {
		msg := "Internal Server Error"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 500, Message: msg})
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *AuthController) GetUserAbilityRules(userType int) []userResponse.Ability {
	if userType == 99 {
		return []userResponse.Ability{{Action: "manage", Subject: "all"}}
	} else {
		return []userResponse.Ability{{Action: "read", Subject: "all"}}
	}
}

func (controller *AuthController) Logout(ctx *gin.Context) {
	// refresh_token, errCookie := ctx.Cookie("refresh_token")

	// if errCookie != nil {
	// 	fileName, atLine := helper.GetFileAndLine(errCookie)
	// 	helper.ResponseError(ctx, helper.CustomError{
	// 		Code:     400,
	// 		Message:  "Invalid Request Structure.",
	// 		FileName: fileName,
	// 		AtLine:   atLine,
	// 	})
	// 	return
	// }

	// env, _ := config.LoadConfig(".")
	// identifier, errValidateToken := helper.ExtractIdentifierFromToken(refresh_token, env.RefreshTokenPublicKey) // you will get UserID & TokenUUID
	// if errValidateToken != nil {
	// 	fileName, atLine := helper.GetFileAndLine(errValidateToken)
	// 	helper.ResponseError(ctx, helper.CustomError{
	// 		Code:     401,
	// 		Message:  "Invalid Token.",
	// 		FileName: fileName,
	// 		AtLine:   atLine,
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    nil,
	})
}

// func (controller *AuthController) VerifyForgetPassword(ctx *gin.Context) {
// 	var payload *authentication.VerifyForgetPassword

// 	errBindJSON := ctx.ShouldBindJSON(&payload)

// 	if errBindJSON != nil {
// 		fileName, atLine := helper.GetFileAndLine(errBindJSON)
// 		helper.ResponseError(ctx, helper.CustomError{
// 			Code:     400,
// 			Message:  "Invalid Request Structure.",
// 			FileName: fileName,
// 			AtLine:   atLine,
// 		})
// 		return
// 	}

// 	errors := helper.ValidateStruct(payload)

// 	if errors != nil {
// 		helper.ResponseError(ctx, helper.CustomError{
// 			Code:     400,
// 			Message:  "Invalid Request Structure.",
// 			FileName: "Auth Controller",
// 			AtLine:   241,
// 		})
// 		return
// 	}

// 	isEmailRegistered := controller.authService.CheckRegisteredEmail(*payload)

// 	if isEmailRegistered {
// 		// Send Verification Email Here
// 		ctx.JSON(http.StatusOK, response.Response{
// 			Success: true,
// 			Code:    200,
// 			Message: "Success",
// 			Data: userResponse.VerifyForgetPassword{
// 				Registered: true,
// 			},
// 		})
// 	} else {
// 		ctx.JSON(http.StatusOK, response.Response{
// 			Success: true,
// 			Code:    200,
// 			Message: "Success",
// 			Data: userResponse.VerifyForgetPassword{
// 				Registered: false,
// 			},
// 		})
// 	}

// }

// func (controller *AuthController) ResetPassword(ctx *gin.Context) {
// 	var payload *authentication.ResetPassword

// 	errBindJSON := ctx.ShouldBindJSON(&payload)

// 	if errBindJSON != nil {
// 		fileName, atLine := helper.GetFileAndLine(errBindJSON)
// 		helper.ResponseError(ctx, helper.CustomError{
// 			Code:     400,
// 			Message:  "Invalid Request Structure.",
// 			FileName: fileName,
// 			AtLine:   atLine,
// 		})
// 		return
// 	}

// 	errors := helper.ValidateStruct(payload)

// 	if errors != nil {
// 		helper.ResponseError(ctx, helper.CustomError{
// 			Code:     400,
// 			Message:  "Invalid Request Structure.",
// 			FileName: "Auth Controller",
// 			AtLine:   241,
// 		})
// 		return
// 	}

// 	errorDB := controller.authService.ResetPassword(*payload)

// 	if errorDB != nil {
// 		helper.ResponseError(ctx, *errorDB)
// 	} else {
// 		ctx.JSON(http.StatusOK, response.Response{
// 			Success: true,
// 			Code:    200,
// 			Message: "Success",
// 			Data: userResponse.ResetPassword{
// 				PasswordValid: true,
// 			},
// 		})
// 	}

// }
func getPositionName(position *model.Position) string {
	if position == nil {
		return ""
	}
	return position.Name
}
