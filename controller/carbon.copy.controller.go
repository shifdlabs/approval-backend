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
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.carbonCopyService.Create(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (controller *CarbonCopyController) Update(ctx *gin.Context) {
	var payload request.CarbonCopyRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.carbonCopyService.Update(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}
