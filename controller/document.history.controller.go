package controller

import (
	documentHistory "Microservice/data/response/DocumentHistory"
	"Microservice/helper"
	service "Microservice/service/DocumentHistory"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHistoryController struct {
	documentHistoryService service.DocumentHistoryService
}

func NewDocumentHistoryController(service service.DocumentHistoryService) *DocumentHistoryController {
	return &DocumentHistoryController{documentHistoryService: service}
}

func (controller *DocumentHistoryController) Get(ctx *gin.Context) {
	stringId := ctx.Param("id")

	documentHistoryResponse, errDocumentHistoryResponse := controller.documentHistoryService.Get(stringId)

	if errDocumentHistoryResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentHistoryResponse)
	} else {
		utils.SuccessResponse(ctx, documentHistoryResponse)
	}
}

func (controller *DocumentHistoryController) GetAll(ctx *gin.Context) {
	cacheData := utils.GetCache(ctx, "All History", &[]documentHistory.DocumentHistoryResponse{})
	if cacheData != nil {
		utils.SuccessResponse(ctx, cacheData)
		return
	}

	documentHistoryResponse, errDocumentHistoryResponse := controller.documentHistoryService.GetAll()

	if errDocumentHistoryResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentHistoryResponse)
	} else {
		utils.SuccessResponse(ctx, documentHistoryResponse)

		utils.SetCache(ctx, "All History", documentHistoryResponse)
	}
}

func (controller *DocumentHistoryController) GetRejectedWithDocumentAndUser(ctx *gin.Context) {
	// Ambil userID dari context
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		msg := "Invalid Request Structure."
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	// Panggil service untuk mendapatkan data
	documentHistoryResponse, err := controller.documentHistoryService.FetchHistoriesByUserID(*id)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	// Kirimkan response sukses
	utils.SuccessResponse(ctx, documentHistoryResponse)
}
