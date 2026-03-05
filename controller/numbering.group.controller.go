package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"
	"fmt"
	"time"

	request "Microservice/data/request/NumberingGroup"
	service "Microservice/service/NumberingGroup"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type NumberingGroupController struct {
	numberingGroupService service.NumberingGroupService
	userLogService        userLogService.UserLogService
}

func NewNumberingGroupController(service service.NumberingGroupService, userLogService userLogService.UserLogService) *NumberingGroupController {
	return &NumberingGroupController{numberingGroupService: service, userLogService: userLogService}
}

func (controller *NumberingGroupController) Get(ctx *gin.Context) {
	stringID := ctx.Param("id")

	numberingGroupResponse, err := controller.numberingGroupService.Get(stringID)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, numberingGroupResponse)
	}
}

func (controller *NumberingGroupController) GetAll(ctx *gin.Context) {
	startTime := time.Now()
	numberingGroupResponse, err := controller.numberingGroupService.GetAll()

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, numberingGroupResponse)
	}
	duration := time.Since(startTime)
	fmt.Printf("GetAll took %s\n", duration)
}

func (controller *NumberingGroupController) Create(ctx *gin.Context) {
	var payload request.NumberingGroupRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.numberingGroupService.Create(payload)

	res := *helper.GetUserUUID(ctx)
	fmt.Println("ID: ", res)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.NumberingGroup),
			Log:    helper.ToJSON(payload),
		},
	)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *NumberingGroupController) Delete(ctx *gin.Context) {
	stringID := ctx.Param("id")
	errResponse := controller.numberingGroupService.Delete(stringID)

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
