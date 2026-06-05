package controller

import (
	"Microservice/helper"
	"Microservice/utils"

	request "Microservice/data/request/AppSettings"
	service "Microservice/service/AppSettings"

	"github.com/gin-gonic/gin"
)

type AppSettingsController struct {
	appSettingsService service.AppSettingService
}

func NewAppSettingsController(service service.AppSettingService) *AppSettingsController {
	return &AppSettingsController{appSettingsService: service}
}

func (controller *AppSettingsController) GetAll(ctx *gin.Context) {
	appSettingsResponse, err := controller.appSettingsService.GetAll()
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, appSettingsResponse)
}

func (controller *AppSettingsController) Update(ctx *gin.Context) {
	var payload request.AppSettingRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.appSettingsService.Update(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}
