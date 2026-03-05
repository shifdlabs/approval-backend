package userlog

import (
	response "Microservice/data/response/UserLog"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/UserLog"

	"github.com/go-playground/validator/v10"
)

type UserLogServiceImpl struct {
	UserLogRepository repository.UserLogRepository
	Validate          *validator.Validate
}

func NewUserLogServiceImpl(
	documentRepository repository.UserLogRepository,
	validate *validator.Validate) UserLogService {
	return &UserLogServiceImpl{
		UserLogRepository: documentRepository,
		Validate:          validate,
	}
}

func (t UserLogServiceImpl) GetAll() ([]response.UserLogResponse, *helper.ErrorModel) {
	result, fetchError := t.UserLogRepository.GetAll()
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapUserLogToUserLogResponse(result), nil
	}
}

func (t UserLogServiceImpl) CreateLog(log model.UserLog) {
	t.UserLogRepository.Create(log)
}
