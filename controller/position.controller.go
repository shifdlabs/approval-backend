package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"
	"fmt"
	"time"

	request "Microservice/data/request/Position"
	service "Microservice/service/Position"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type PositionController struct {
	positionService service.PositionService
	userLogService  userLogService.UserLogService
}

func NewPositionController(service service.PositionService, userLogService userLogService.UserLogService) *PositionController {
	return &PositionController{positionService: service, userLogService: userLogService}
}

func (controller *PositionController) Get(ctx *gin.Context) {
	stringID := ctx.Param("id")

	positionResponse, err := controller.positionService.Get(stringID)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, positionResponse)
	}
}

func (controller *PositionController) GetAll(ctx *gin.Context) {
	startTime := time.Now()
	positionResponse, err := controller.positionService.GetAll()

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, positionResponse)
	}
	duration := time.Since(startTime)
	fmt.Printf("GetAll took %s\n", duration)
}

func (controller *PositionController) Create(ctx *gin.Context) {
	var payload request.CreatePositionRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
	}

	err := controller.positionService.Create(payload)

	res := *helper.GetUserUUID(ctx)
	fmt.Println("ID: ", res)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.Position),
			Log:    helper.ToJSON(payload),
		},
	)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *PositionController) Update(ctx *gin.Context) {
	var payload request.UpdatePositionRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.positionService.Update(payload)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Update),
			Module: string(enums.Position),
			Log:    helper.ToJSON(payload),
		},
	)

	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

func (controller *PositionController) Delete(ctx *gin.Context) {
	stringID := ctx.Param("id")
	errResponse := controller.positionService.Delete(stringID)

	// Action Log
	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.Position),
		},
	)

	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}
