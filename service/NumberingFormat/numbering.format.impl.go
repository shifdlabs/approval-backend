package numberingformat

import (
	request "Microservice/data/request/NumberingFormat"
	response "Microservice/data/response/NumberingFormat"
	"Microservice/helper"
)

type NumberingFormatService interface {
	Create(request request.NumberingFormatRequest) *helper.ErrorModel
	GetAll() ([]response.NumberingFormatResponse, *helper.ErrorModel)
	GetAllWithGrouped() ([]response.NumberingFormatByGroupResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
