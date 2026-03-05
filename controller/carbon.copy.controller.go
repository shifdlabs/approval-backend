package controller

import (
	"Microservice/helper"
	"Microservice/utils"

	request "Microservice/data/request/CarbonCopy"
	service "Microservice/service/CarbonCopy"

	"github.com/gin-gonic/gin"
)

type CarbonCopyController struct {
	carbonCopyService service.CarbonCopyService
}

func NewCarbonCopyController(service service.CarbonCopyService) *CarbonCopyController {
	return &CarbonCopyController{carbonCopyService: service}
}

func (controller *CarbonCopyController) Create(ctx *gin.Context) {
	var payload request.CarbonCopyRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.carbonCopyService.Create(payload)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *CarbonCopyController) Update(ctx *gin.Context) {
	var payload request.CarbonCopyRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.carbonCopyService.Update(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
