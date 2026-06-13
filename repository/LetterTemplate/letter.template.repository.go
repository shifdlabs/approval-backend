package lettertemplate

import (
	"Microservice/helper"
	"Microservice/model"
)

type LetterTemplateRepository interface {
	GetAll() ([]model.LetterTemplate, *helper.ErrorModel)
	GetByID(id string) (*model.LetterTemplate, *helper.ErrorModel)
	Create(template model.LetterTemplate) (*model.LetterTemplate, *helper.ErrorModel)
	Update(id string, template model.LetterTemplate) (*model.LetterTemplate, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
