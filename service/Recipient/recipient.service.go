package recipient

import (
	request "Microservice/data/request/Recipient"
	"Microservice/helper"
)

type RecipientService interface {
	Create(request request.RecipientRequest) *helper.ErrorModel
	Update(request request.RecipientRequest) *helper.ErrorModel
}
