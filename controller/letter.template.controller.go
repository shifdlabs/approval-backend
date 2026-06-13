package controller

import (
	request "Microservice/data/request/LetterTemplate"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/LetterTemplate"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type LetterTemplateController struct {
	service        service.LetterTemplateService
	userLogService userLogService.UserLogService
}

func NewLetterTemplateController(svc service.LetterTemplateService, logSvc userLogService.UserLogService) *LetterTemplateController {
	return &LetterTemplateController{service: svc, userLogService: logSvc}
}

func (c *LetterTemplateController) GetAll(ctx *gin.Context) {
	templates, err := c.service.GetAll()
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, templates)
}

func (c *LetterTemplateController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	template, err := c.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, template)
}

func (c *LetterTemplateController) Create(ctx *gin.Context) {
	var payload request.CreateLetterTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	result, err := c.service.Create(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Create),
		Module: "Letter Template",
	})
	utils.SuccessResponse(ctx, result)
}

func (c *LetterTemplateController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var payload request.UpdateLetterTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	result, err := c.service.Update(id, payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: "Letter Template",
	})
	utils.SuccessResponse(ctx, result)
}

func (c *LetterTemplateController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Delete),
		Module: "Letter Template",
	})
	utils.SuccessResponse(ctx, nil)
}
