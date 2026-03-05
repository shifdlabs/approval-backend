package userlog

import (
	response "Microservice/data/response/UserLog"
	"Microservice/helper"
	"Microservice/model"
)

type UserLogService interface {
	GetAll() ([]response.UserLogResponse, *helper.ErrorModel)
	CreateLog(log model.UserLog)
}
