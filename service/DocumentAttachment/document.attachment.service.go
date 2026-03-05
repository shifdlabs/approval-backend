package documentAttachment

import (
	response "Microservice/data/response/DocumentAttachment"
	"Microservice/helper"
)

type DocumentAttachmentService interface {
	Get(id string) (*response.DocumentAttachmentResponse, *helper.ErrorModel)
	GetAll() ([]response.DocumentAttachmentResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
}
