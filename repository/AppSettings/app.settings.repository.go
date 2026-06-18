package appsettings

import (
	"Microservice/helper"
	"Microservice/model"
)

type AppSettingsRepository interface {
	GetAll() ([]model.AppSettings, *helper.ErrorModel)
	GetByKey(key string) (*model.AppSettings, *helper.ErrorModel)
	Update(properties []model.AppSettings) *helper.ErrorModel
}
