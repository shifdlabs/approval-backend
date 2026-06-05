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
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	err := controller.documentNumbersService.Create(payload, *id, nil, enums.Booked)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.DocumentNumbers),
			Log:    helper.ToJSON(payload),
		},
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *DocumentNumbersController) GetAll(ctx *gin.Context) {
	documentNumbers, errDocumentNumbersResponse := controller.documentNumbersService.GetAll()
	if errDocumentNumbersResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentNumbersResponse)
		return
	}
	utils.SuccessResponse(ctx, documentNumbers)
}

func (controller *DocumentNumbersController) GetAllByUserId(ctx *gin.Context) {
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentNumbers, errDocumentNumbersResponse := controller.documentNumbersService.GetAllByUserId(*id)
	if errDocumentNumbersResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentNumbersResponse)
		return
	}
	utils.SuccessResponse(ctx, documentNumbers)
}

func (controller *DocumentNumbersController) Delete(ctx *gin.Context) {
	stringID := ctx.Param("id")
	errResponse := controller.documentNumbersService.Delete(stringID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.DocumentNumbers),
		},
	)

	utils.SuccessResponse(ctx, nil)
}
