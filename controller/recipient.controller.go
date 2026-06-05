package controller

import (
	"Microservice/helper"
	"Microservice/utils"

	request "Microservice/data/request/Recipient"
	service "Microservice/service/Recipient"

	"github.com/gin-gonic/gin"
)

type RecipientController struct {
	recipientService service.RecipientService
}

func NewRecipientController(service service.RecipientService) *RecipientController {
	return &RecipientController{recipientService: service}
}

func (controller *RecipientController) Create(ctx *gin.Context) {
	var payload request.RecipientRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.recipientService.Create(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (controller *RecipientController) Update(ctx *gin.Context) {
	var payload request.RecipientRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.recipientService.Update(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}
