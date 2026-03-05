package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"

	request "Microservice/data/request/DocumentNumbers"
	service "Microservice/service/DocumentNumbers"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type DocumentNumbersController struct {
	documentNumbersService service.DocumentNumbersService
	userLogService         userLogService.UserLogService
}

func NewDocumentNumbersController(service service.DocumentNumbersService, userLogService userLogService.UserLogService) *DocumentNumbersController {
	return &DocumentNumbersController{documentNumbersService: service, userLogService: userLogService}
}

func (controller *DocumentNumbersController) Create(ctx *gin.Context) {
	var payload request.DocumentNumbersRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		msg := "Invalid Request Structure."
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.documentNumbersService.Create(payload, *id, nil, enums.Booked)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.DocumentNumbers),
			Log:    helper.ToJSON(payload),
		},
	)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *DocumentNumbersController) GetAll(ctx *gin.Context) {

	documentNumbers, errDocumentNumbersResponse := controller.documentNumbersService.GetAll()

	if errDocumentNumbersResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentNumbersResponse)
	} else {
		utils.SuccessResponse(ctx, documentNumbers)
	}
}

func (controller *DocumentNumbersController) GetAllByUserId(ctx *gin.Context) {
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		msg := "Invalid Request Structure."
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	documentNumbers, errDocumentNumbersResponse := controller.documentNumbersService.GetAllByUserId(*id)

	if errDocumentNumbersResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentNumbersResponse)
	} else {
		utils.SuccessResponse(ctx, documentNumbers)
	}
}

func (controller *DocumentNumbersController) Delete(ctx *gin.Context) {
	stringID := ctx.Param("id")
	errResponse := controller.documentNumbersService.Delete(stringID)
	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.DocumentNumbers),
		},
	)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
