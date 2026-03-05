package controller

import (
	request "Microservice/data/request/Attachment"
	documentAttachment "Microservice/data/response/DocumentAttachment"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/DocumentAttachment"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentAttachmentController struct {
	documentAttachmentService service.DocumentAttachmentService
	userLogService            userLogService.UserLogService
}

func NewDocumentAttachmentController(service service.DocumentAttachmentService, userLogService userLogService.UserLogService) *DocumentAttachmentController {
	return &DocumentAttachmentController{documentAttachmentService: service, userLogService: userLogService}
}

func (controller *DocumentAttachmentController) Get(ctx *gin.Context) {
	stringId := ctx.Param("id")

	documentAttachmentResponse, errDocumentAttachmentResponse := controller.documentAttachmentService.Get(stringId)

	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
	} else {
		utils.SuccessResponse(ctx, documentAttachmentResponse)
	}
}

func (controller *DocumentAttachmentController) GetAll(ctx *gin.Context) {
	cacheData := utils.GetCache(ctx, "All Attachment", &[]documentAttachment.DocumentAttachmentResponse{})
	if cacheData != nil {
		utils.SuccessResponse(ctx, cacheData)
		return
	}

	documentAttachmentResponse, errDocumentAttachmentResponse := controller.documentAttachmentService.GetAll()

	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
	} else {
		utils.SuccessResponse(ctx, documentAttachmentResponse)

		utils.SetCache(ctx, "All Attachment", documentAttachmentResponse)
	}
}

func (controller *DocumentAttachmentController) Delete(ctx *gin.Context) {
	var payload request.AttachmentRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	errDocumentAttachmentResponse := controller.documentAttachmentService.Delete(payload.Id)

	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
	} else {
		controller.userLogService.CreateLog(
			model.UserLog{
				UserID: *helper.GetUserUUID(ctx),
				Action: string(enums.Delete),
				Module: string(enums.DocumentAttachment),
				Log:    helper.ToJSON(payload),
			},
		)

		utils.SuccessResponse(ctx, nil)
	}
}
