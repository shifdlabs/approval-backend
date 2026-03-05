package appsettings

import (
	request "Microservice/data/request/AppSettings"
	response "Microservice/data/response/AppSettings"
	"Microservice/helper"
)

type AppSettingService interface {
	GetAll() ([]response.AppSettingResponse, *helper.ErrorModel)
	Update(appSettings request.AppSettingRequest) *helper.ErrorModel
}
