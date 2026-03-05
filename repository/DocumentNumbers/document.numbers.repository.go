package documentnumbers

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersRepository interface {
	Create(data model.DocumentNumbers) *helper.ErrorModel
	Get(id string) (*model.DocumentNumbers, *helper.ErrorModel)
	GetByDocumentID(id uuid.UUID) (*model.DocumentNumbers, *helper.ErrorModel)
	GetAll() ([]model.DocumentNumbers, *helper.ErrorModel)
	GetAllByUserID(userId string) ([]model.DocumentNumbers, *helper.ErrorModel)
	GetTotal(formatId string, groupId *string) (*int64, *helper.ErrorModel)
	GetCancelled(formatId string, groupId *string) (*model.DocumentNumbers, *helper.ErrorModel)
	Update(data model.DocumentNumbers) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
}
