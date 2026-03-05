package signature

import (
	"Microservice/helper"
	"Microservice/model"
)

type SignatureRepository interface {
	Create(signature *model.Signature) *helper.ErrorModel
	Update(signature *model.Signature) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	GetAll() ([]model.Signature, *helper.ErrorModel)
	GetByUserId(userId string) (*model.Signature, *helper.ErrorModel)
	GetByUserIds(userIds []string) ([]model.Signature, *helper.ErrorModel)
}
