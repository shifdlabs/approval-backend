package controller

import (
	request "Microservice/data/request/Delegator"
	"Microservice/helper"
	service "Microservice/service/Delegator"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DelegatorController struct {
	delegatorService service.DelegatorService
}

func NewDelegatorController(svc service.DelegatorService) *DelegatorController {
	return &DelegatorController{delegatorService: svc}
}

func (c *DelegatorController) GetAll(ctx *gin.Context) {
	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	result, err := c.delegatorService.GetAll(*ownerID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, result)
}

func (c *DelegatorController) Create(ctx *gin.Context) {
	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	var payload request.CreateDelegatorRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := c.delegatorService.Create(*ownerID, payload); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (c *DelegatorController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	var payload request.UpdateDelegatorRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := c.delegatorService.Update(id, *ownerID, payload); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (c *DelegatorController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if err := c.delegatorService.Delete(id, *ownerID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}
