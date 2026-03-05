package controller

import (
	signatureRequest "Microservice/data/request/Signature"
	"Microservice/helper"
	signatureService "Microservice/service/Signature"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type SignatureController struct {
	signatureService signatureService.SignatureService
}

func NewSignatureController(service signatureService.SignatureService) *SignatureController {
	return &SignatureController{signatureService: service}
}

func (controller *SignatureController) Create(ctx *gin.Context) {
	var payload signatureRequest.CreateSignatureRequest

	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.signatureService.Create(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) Update(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var payload signatureRequest.UpdateSignatureRequest

	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.signatureService.Update(userId, payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) Delete(ctx *gin.Context) {
	userId := ctx.Param("userId")

	err := controller.signatureService.Delete(userId)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) GetAll(ctx *gin.Context) {
	signatures, err := controller.signatureService.GetAll()
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, signatures)
}

func (controller *SignatureController) GetByUserId(ctx *gin.Context) {
	userId := ctx.Param("userId")

	signature, err := controller.signatureService.GetByUserId(userId)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	if signature == nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 404, Message: "Signature not found"})
		return
	}

	utils.SuccessResponse(ctx, signature)
}
