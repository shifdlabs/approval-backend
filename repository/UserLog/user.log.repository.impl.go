package userlog

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type UserLogRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserLogRepositoryImpl(Db *gorm.DB) UserLogRepository {
	return &UserLogRepositoryImpl{Db: Db}
}

func (t *UserLogRepositoryImpl) Create(document model.UserLog) {
	result := t.Db.Create(&document)

	if result.Error != nil {
		msg := "Failed to create user log"
		helper.ErrorLog(result.Error, 500, &msg)
	}
}

func (t *UserLogRepositoryImpl) GetAll() ([]model.UserLog, *helper.ErrorModel) {
	var userLogs []model.UserLog
	result := t.Db.Find(&userLogs)
	if result.Error != nil {
		msg := "Failed to get all user logs"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return userLogs, nil
}
