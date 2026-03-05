package numberingformat

import (
	"Microservice/helper"
	"Microservice/model"
)

type NumberingFormatRepository interface {
	Create(data model.NumberingFormat) *helper.ErrorModel
	Get(id string) (*model.NumberingFormat, *helper.ErrorModel)
	GetAll() ([]model.NumberingFormat, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
