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
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.recipientService.Create(payload)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *RecipientController) Update(ctx *gin.Context) {
	var payload request.RecipientRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.recipientService.Update(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
