package position

import (
	"Microservice/helper"
	"Microservice/model"
)

type PositionRepository interface {
	Create(report model.Position) *helper.ErrorModel
	Get(id string) (*model.Position, *helper.ErrorModel)
	GetAll() ([]model.Position, *helper.ErrorModel)
	Update(position model.Position) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
}
