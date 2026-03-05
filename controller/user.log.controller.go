package controller

import (
	service "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type UserLogController struct {
	userLogService service.UserLogService
}

func NewUserLogController(service service.UserLogService) *UserLogController {
	return &UserLogController{userLogService: service}
}

func (controller *UserLogController) GetAll(ctx *gin.Context) {
	userLogResponse, errUserLogResponse := controller.userLogService.GetAll()

	if errUserLogResponse != nil {
		utils.ErrorResponse(ctx, *errUserLogResponse)
	} else {
		utils.SuccessResponse(ctx, userLogResponse)
	}
}
