package delegator

import (
	request "Microservice/data/request/Delegator"
	response "Microservice/data/response/Delegator"
	"Microservice/helper"
	"time"
)

type DelegatorService interface {
	Create(ownerID string, req request.CreateDelegatorRequest) *helper.ErrorModel
	GetAll(ownerID string) ([]response.DelegatorResponse, *helper.ErrorModel)
	Update(id string, ownerID string, req request.UpdateDelegatorRequest) *helper.ErrorModel
	Delete(id string, ownerID string) *helper.ErrorModel
	ResolveDelegate(ownerID string, date time.Time) (string, *helper.ErrorModel)
	GetOwnerChainForDelegate(delegateID string, date time.Time) ([]string, *helper.ErrorModel)
}
