package controller

import (
	request "Microservice/data/request/User"
	"Microservice/helper"
	service "Microservice/service/User"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{userService: service}
}

func (controller *UserController) Get(ctx *gin.Context) {
	// GET USER ID
	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	userResponse, errResponse := controller.userService.Get(*stringID)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, userResponse)
	}
}

func (controller *UserController) GetUserByID(ctx *gin.Context) {
	// GET USER ID
	stringID := ctx.Param("id")
	if stringID == "" {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	userResponse, errResponse := controller.userService.Get(stringID)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, userResponse)
	}
}

func (controller *UserController) GetAll(ctx *gin.Context) {
	userResponse, errResponse := controller.userService.GetAll()

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, userResponse)
	}
}

func (controller *UserController) GetAllUserExceptCurrent(ctx *gin.Context) {
	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	userResponse, errResponse := controller.userService.GetAllUserExceptCurrent(*stringID)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, userResponse)
	}
}

func (controller *UserController) Create(ctx *gin.Context) {
	var payload request.CreateUserRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.Create(payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) Update(ctx *gin.Context) {
	var payload request.UpdateUserRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.Update(payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) Delete(ctx *gin.Context) {
	stringId := ctx.Param("id")

	err := controller.userService.Delete(stringId)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) MultipleDelete(ctx *gin.Context) {
	var payload request.DeleteMultipleUserRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.userService.MultipleDelete(payload.IDs)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) UpdateEmail(ctx *gin.Context) {
	var payload request.UpdateEmailRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.UpdateEmail(*stringID, payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) UpdateBiodata(ctx *gin.Context) {
	var payload request.UpdateBiodataRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.UpdateBiodata(*stringID, payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) UpdateRole(ctx *gin.Context) {
	var payload request.UpdateRoleRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.UpdateRole(payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) UpdatePassword(ctx *gin.Context) {
	var payload request.UpdatePasswordRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.UpdatePassword(payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *UserController) UpdateAccess(ctx *gin.Context) {
	var payload request.UpdateAccessRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	errCreateUser := controller.userService.UpdateAccess(payload)
	if errCreateUser != nil {
		utils.ErrorResponse(ctx, *errCreateUser)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
