package numberinggroup

import (
	request "Microservice/data/request/NumberingGroup"
	response "Microservice/data/response/NumberingGroup"
	"Microservice/helper"
)

type NumberingGroupService interface {
	Create(request request.NumberingGroupRequest) *helper.ErrorModel
	Get(id string) (*response.NumberingGroupResponse, *helper.ErrorModel)
	GetAll() ([]response.NumberingGroupResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
