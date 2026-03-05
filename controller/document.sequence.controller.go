package controller

import (
	"Microservice/helper"
	service "Microservice/service/DocumentSequence"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentSequenceController struct {
	documentSequenceService service.DocumentSequenceService
}

func NewDocumentSequenceController(service service.DocumentSequenceService) *DocumentSequenceController {
	return &DocumentSequenceController{documentSequenceService: service}
}

func (controller *DocumentSequenceController) Get(ctx *gin.Context) {
	stringId := ctx.Param("id")

	documentSequenceResponse, errDocumentSequenceResponse := controller.documentSequenceService.Get(stringId)

	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
	} else {
		utils.SuccessResponse(ctx, documentSequenceResponse)
	}
}

func (controller *DocumentSequenceController) GetProgress(ctx *gin.Context) {
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		msg := "Invalid Request Structure."
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	// Panggil service untuk mendapatkan data progress
	documentSequenceResponses, err := controller.documentSequenceService.GetProgressByAuthorID(*id)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	// Kirimkan response sukses
	utils.SuccessResponse(ctx, documentSequenceResponses)
}
