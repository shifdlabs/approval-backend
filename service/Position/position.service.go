package position

import (
	request "Microservice/data/request/Position"
	response "Microservice/data/response/Position"
	"Microservice/helper"
)

type PositionService interface {
	Create(position request.CreatePositionRequest) *helper.ErrorModel
	Get(id string) (*response.PositionResponse, *helper.ErrorModel)
	GetAll() ([]response.PositionResponse, *helper.ErrorModel)
	Update(position request.UpdatePositionRequest) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
}
