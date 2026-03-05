package documentsequence

import (
	"Microservice/helper"
	"Microservice/model"
)

type NumberingGroupRepository interface {
	Create(data model.NumberingGroup) *helper.ErrorModel
	Get(id string) (*model.NumberingGroup, *helper.ErrorModel)
	GetAll() ([]model.NumberingGroup, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
