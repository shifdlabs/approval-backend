package userlog

import (
	"Microservice/helper"
	"Microservice/model"
)

type UserLogRepository interface {
	Create(document model.UserLog)
	GetAll() ([]model.UserLog, *helper.ErrorModel)
}
