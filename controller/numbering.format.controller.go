package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"
	"fmt"

	request "Microservice/data/request/NumberingFormat"
	service "Microservice/service/NumberingFormat"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type NumberingFormatController struct {
	numberingFormatService service.NumberingFormatService
	userLogService         userLogService.UserLogService
}

func NewNumberingFormatController(service service.NumberingFormatService, userLogService userLogService.UserLogService) *NumberingFormatController {
	return &NumberingFormatController{numberingFormatService: service, userLogService: userLogService}
}

func (controller *NumberingFormatController) Create(ctx *gin.Context) {
	var payload request.NumberingFormatRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.numberingFormatService.Create(payload)

	res := *helper.GetUserUUID(ctx)
	fmt.Println("ID: ", res)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.NumberingFormat),
			Log:    helper.ToJSON(payload),
		},
	)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *NumberingFormatController) GetAll(ctx *gin.Context) {
	documentSequenceResponse, errDocumentSequenceResponse := controller.numberingFormatService.GetAll()

	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
	} else {
		utils.SuccessResponse(ctx, documentSequenceResponse)
	}
}

func (controller *NumberingFormatController) GetAllWithGrouped(ctx *gin.Context) {
	documentSequenceResponse, errDocumentSequenceResponse := controller.numberingFormatService.GetAllWithGrouped()

	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
	} else {
		utils.SuccessResponse(ctx, documentSequenceResponse)
	}
}

func (controller *NumberingFormatController) Delete(ctx *gin.Context) {
	stringID := ctx.Param("id")
	errResponse := controller.numberingFormatService.Delete(stringID)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.NumberingGroup),
		},
	)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
