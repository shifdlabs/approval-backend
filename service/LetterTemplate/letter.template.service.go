package lettertemplate

import (
	request "Microservice/data/request/LetterTemplate"
	response "Microservice/data/response/LetterTemplate"
	"Microservice/helper"
)

type LetterTemplateService interface {
	GetAll() ([]response.LetterTemplateResponse, *helper.ErrorModel)
	GetByID(id string) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Create(req request.CreateLetterTemplateRequest) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Update(id string, req request.UpdateLetterTemplateRequest) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
