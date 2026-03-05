package controller

import (
	"Microservice/helper"
	"Microservice/utils"
	"fmt"
	"time"

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
	startTime := time.Now()
	appSettingsResponse, err := controller.appSettingsService.GetAll()

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, appSettingsResponse)
	}
	duration := time.Since(startTime)
	fmt.Printf("GetAll took %s\n", duration)
}

func (controller *AppSettingsController) Update(ctx *gin.Context) {
	var payload request.AppSettingRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.appSettingsService.Update(payload)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
