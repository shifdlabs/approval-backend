package delegator

import (
	"Microservice/helper"
	"Microservice/model"
	"time"
)

type DelegatorRepository interface {
	Create(delegator model.Delegator) *helper.ErrorModel
	GetAllByOwnerID(ownerID string) ([]model.Delegator, *helper.ErrorModel)
	GetByID(id string) (*model.Delegator, *helper.ErrorModel)
	Update(delegator model.Delegator) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	GetActiveDelegationByOwnerID(ownerID string, date time.Time) (*model.Delegator, *helper.ErrorModel)
	GetOwnerIDsByDelegateID(delegateID string, date time.Time) ([]string, *helper.ErrorModel)
	HasOverlappingDelegation(ownerID string, startDate time.Time, endDate time.Time, excludeID *string) (bool, *helper.ErrorModel)
}
