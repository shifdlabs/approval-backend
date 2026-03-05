package signature

import (
	signatureRequest "Microservice/data/request/Signature"
	signatureResponse "Microservice/data/response/Signature"
	"Microservice/helper"
)

type SignatureService interface {
	Create(request signatureRequest.CreateSignatureRequest) *helper.ErrorModel
	Update(userId string, request signatureRequest.UpdateSignatureRequest) *helper.ErrorModel
	Delete(userId string) *helper.ErrorModel
	GetAll() ([]signatureResponse.SignatureResponse, *helper.ErrorModel)
	GetByUserId(userId string) (*signatureResponse.SignatureResponse, *helper.ErrorModel)
}
