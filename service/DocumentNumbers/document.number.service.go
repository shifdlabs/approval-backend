package documentnumbers

import (
	request "Microservice/data/request/DocumentNumbers"
	response "Microservice/data/response/DocumentNumbers"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersService interface {
	Create(request request.DocumentNumbersRequest, userId string, document *model.Document, state enums.DocumentNumberState) *helper.ErrorModel
	Update(id string, document *model.Document, state enums.DocumentNumberState) *helper.ErrorModel
	GetAll() ([]response.DocumentNumbersResponse, *helper.ErrorModel)
	Get(id string) (*response.DocumentNumbersResponse, *helper.ErrorModel)
	GetByDocumentID(id uuid.UUID) (*response.DocumentNumbersResponse, *helper.ErrorModel)
	GetAllByUserId(userId string) ([]response.DocumentNumbersResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
