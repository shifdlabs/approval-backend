package userlog

import (
	"Microservice/helper"
	"Microservice/model"
)

type UserLogWithName struct {
	model.UserLog
	UserName string
}

type UserLogRepository interface {
	Create(document model.UserLog)
	GetAll() ([]UserLogWithName, *helper.ErrorModel)
}
