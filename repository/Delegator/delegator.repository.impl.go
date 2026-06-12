package delegator

import (
	"Microservice/helper"
	"Microservice/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DelegatorRepositoryImpl struct {
	Db *gorm.DB
}

func NewDelegatorRepositoryImpl(db *gorm.DB) DelegatorRepository {
	return &DelegatorRepositoryImpl{Db: db}
}

func (r *DelegatorRepositoryImpl) Create(delegator model.Delegator) *helper.ErrorModel {
	result := r.Db.Create(&delegator)
	if result.Error != nil {
		msg := "Failed to create delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) GetAllByOwnerID(ownerID string) ([]model.Delegator, *helper.ErrorModel) {
	var delegators []model.Delegator
	result := r.Db.Preload("Owner").Preload("Delegate").
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&delegators)
	if result.Error != nil {
		msg := "Failed to get delegations"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return delegators, nil
}

func (r *DelegatorRepositoryImpl) GetByID(id string) (*model.Delegator, *helper.ErrorModel) {
	var delegator model.Delegator
	delegatorID, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to parse UUID"
		return nil, helper.ErrorCatcher(errParse, 400, &msg)
	}
	result := r.Db.Preload("Owner").Preload("Delegate").First(&delegator, "id = ?", delegatorID)
	if result.Error != nil {
		msg := "Delegation not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}
	return &delegator, nil
}

func (r *DelegatorRepositoryImpl) Update(delegator model.Delegator) *helper.ErrorModel {
	result := r.Db.Model(&delegator).Updates(map[string]interface{}{
		"delegate_id": delegator.DelegateID,
		"start_date":  delegator.StartDate,
		"end_date":    delegator.EndDate,
	})
	if result.Error != nil {
		msg := "Failed to update delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) Delete(id string) *helper.ErrorModel {
	delegatorID, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to parse UUID"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}
	result := r.Db.Unscoped().Delete(&model.Delegator{}, delegatorID)
	if result.Error != nil {
		msg := "Failed to delete delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) GetActiveDelegationByOwnerID(ownerID string, date time.Time) (*model.Delegator, *helper.ErrorModel) {
	var delegator model.Delegator
	result := r.Db.Preload("Delegate").
		Where("owner_id = ? AND start_date <= ? AND end_date >= ?", ownerID, date, date).
		First(&delegator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		msg := "Failed to get active delegation"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return &delegator, nil
}

func (r *DelegatorRepositoryImpl) GetOwnerIDsByDelegateID(delegateID string, date time.Time) ([]string, *helper.ErrorModel) {
	var ownerIDs []string
	result := r.Db.Model(&model.Delegator{}).
		Select("owner_id").
		Where("delegate_id = ? AND start_date <= ? AND end_date >= ?", delegateID, date, date).
		Pluck("owner_id", &ownerIDs)
	if result.Error != nil {
		msg := "Failed to get owners by delegate"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return ownerIDs, nil
}

func (r *DelegatorRepositoryImpl) HasOverlappingDelegation(ownerID string, startDate time.Time, endDate time.Time, excludeID *string) (bool, *helper.ErrorModel) {
	query := r.Db.Model(&model.Delegator{}).
		Where("owner_id = ? AND start_date <= ? AND end_date >= ?", ownerID, endDate, startDate)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		msg := "Failed to check overlapping delegation"
		return false, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return count > 0, nil
}
