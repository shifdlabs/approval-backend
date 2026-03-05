package documentSequence

import (
	response "Microservice/data/response/DocumentSequence"
	"Microservice/helper"
)

type DocumentSequenceService interface {
	Get(id string) (*response.DocumentSequenceResponse, *helper.ErrorModel)
	GetAll() ([]response.DocumentSequenceResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
	GetProgressByAuthorID(authorID string) ([]response.DocumentSequenceResponse, *helper.ErrorModel)
}
