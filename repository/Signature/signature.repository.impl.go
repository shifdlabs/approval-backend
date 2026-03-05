package signature

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type SignatureRepositoryImpl struct {
	Db *gorm.DB
}

func NewSignatureRepositoryImpl(db *gorm.DB) SignatureRepository {
	return &SignatureRepositoryImpl{Db: db}
}

func (t *SignatureRepositoryImpl) Create(signature *model.Signature) *helper.ErrorModel {
	result := t.Db.Create(signature)
	if result.Error != nil {
		msg := "Failed to create signature"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (t *SignatureRepositoryImpl) Update(signature *model.Signature) *helper.ErrorModel {
	result := t.Db.Save(signature)
	if result.Error != nil {
		msg := "Failed to update signature"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (t *SignatureRepositoryImpl) Delete(id string) *helper.ErrorModel {
	signatureId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Delete(&model.Signature{}, signatureId)
	if result.Error != nil {
		msg := "Failed to delete signature"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (t *SignatureRepositoryImpl) GetByUserId(userId string) (*model.Signature, *helper.ErrorModel) {
	var signature model.Signature
	userIdUUID, err := uuid.FromString(userId)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Where("user_id = ?", userIdUUID).First(&signature)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		msg := "Signature not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &signature, nil
}

func (t *SignatureRepositoryImpl) GetAll() ([]model.Signature, *helper.ErrorModel) {
	var signatures []model.Signature
	result := t.Db.Find(&signatures)
	if result.Error != nil {
		msg := "Failed to get signatures"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return signatures, nil
}

func (t *SignatureRepositoryImpl) GetByUserIds(userIds []string) ([]model.Signature, *helper.ErrorModel) {
	var signatures []model.Signature
	var userUUIDs []uuid.UUID

	for _, userId := range userIds {
		userIdUUID, err := uuid.FromString(userId)
		if err != nil {
			continue
		}
		userUUIDs = append(userUUIDs, userIdUUID)
	}

	result := t.Db.Where("user_id IN ?", userUUIDs).Find(&signatures)
	if result.Error != nil {
		msg := "Failed to get signatures"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return signatures, nil
}
