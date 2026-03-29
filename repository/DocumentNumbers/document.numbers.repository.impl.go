package documentnumbers

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentNumbersRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentNumbersRepositoryImpl(Db *gorm.DB) DocumentNumbersRepository {
	return &DocumentNumbersRepositoryImpl{Db: Db}
}

func (t *DocumentNumbersRepositoryImpl) Create(data model.DocumentNumbers) *helper.ErrorModel {
	result := t.Db.Save(&data)

	if result.Error != nil {
		msg := "Failed to create numbering format"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentNumbersRepositoryImpl) Update(data model.DocumentNumbers) *helper.ErrorModel {
	var existing model.DocumentNumbers
	if err := t.Db.First(&existing, data.ID).Error; err != nil {
		msg := "User not found"
		return helper.ErrorCatcher(err, 404, &msg)
	}

	// We have to add .Select("*") so gorm will not ignoring zero value like 'false', and it can still updating all value
	result := t.Db.Model(&existing).Select("*").Updates(data)
	if result.Error != nil {
		msg := "Failed to Update User Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentNumbersRepositoryImpl) Get(id string) (*model.DocumentNumbers, *helper.ErrorModel) {
	var documentNumbers model.DocumentNumbers
	documentNumbersId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.First(&documentNumbers, "id = ?", documentNumbersId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Document Number ID Not Found"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentNumbers, nil
}

func (t *DocumentNumbersRepositoryImpl) GetByDocumentID(id uuid.UUID) (*model.DocumentNumbers, *helper.ErrorModel) {
	var documentNumbers model.DocumentNumbers
	result := t.Db.Preload("NumberingFormat").First(&documentNumbers, "document_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Document Number ID Not Found"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentNumbers, nil
}

func (t *DocumentNumbersRepositoryImpl) GetAll() ([]model.DocumentNumbers, *helper.ErrorModel) {
	var documentNumberss []model.DocumentNumbers
	result := t.Db.Where("deleted_at IS NULL").Find(&documentNumberss)
	if result.Error != nil {
		msg := "Failed to get all numbering formats"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentNumberss, nil
}

func (t *DocumentNumbersRepositoryImpl) GetTotal(formatId string, groupId *string) (*int64, *helper.ErrorModel) {
	var totalRecord int64
	if groupId != nil {
		countTotalRecord := t.Db.Model(&model.DocumentNumbers{}).
			Joins("JOIN numbering_formats ON numbering_formats.id = document_numbers.numbering_format_id").
			Where("numbering_formats.group_id = ?", groupId).
			Where("numbering_formats.increment_by_group = ?", true).
			Count(&totalRecord)

		if countTotalRecord.Error != nil {
			msg := "Failed to get all documents"
			return nil, helper.ErrorCatcher(countTotalRecord.Error, 500, &msg)
		}

		return &totalRecord, nil
	} else {
		countTotalRecord := t.Db.Model(&model.DocumentNumbers{}).
			Where("numbering_format_id = ?", formatId).
			Count(&totalRecord)

		if countTotalRecord.Error != nil {
			msg := "Failed to get all documents"
			return nil, helper.ErrorCatcher(countTotalRecord.Error, 500, &msg)
		}

		return &totalRecord, nil
	}
}

func (t *DocumentNumbersRepositoryImpl) GetCancelled(formatId string, groupId *string) (*model.DocumentNumbers, *helper.ErrorModel) {
	documentNumber := &model.DocumentNumbers{}

	if groupId != nil {
		cancelledRecord := t.Db.Model(&model.DocumentNumbers{}).
			Joins("JOIN numbering_formats ON numbering_formats.id = document_numbers.numbering_format_id").
			Where("numbering_formats.group_id = ?", groupId).
			Where("numbering_formats.id = ?", formatId).
			Where("document_numbers.state = 0").
			Where("numbering_formats.increment_by_group = ?", true).
			First(documentNumber)

		if cancelledRecord.Error != nil {
			if errors.Is(cancelledRecord.Error, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			msg := "Document Number not found"
			return nil, helper.ErrorCatcher(cancelledRecord.Error, 404, &msg)
		}

		return documentNumber, nil
	} else {
		formatUUID, err := uuid.FromString(formatId)
		if err != nil {
			msg := "Failed to parse uuid"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}

		cancelledRecord := t.Db.Model(&model.DocumentNumbers{}).
			Where("numbering_format_id = ?", formatUUID).
			Where("document_numbers.state = 0").
			First(documentNumber)

		if cancelledRecord.Error != nil {
			if errors.Is(cancelledRecord.Error, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			msg := "Document Number not found"
			return nil, helper.ErrorCatcher(cancelledRecord.Error, 404, &msg)
		}

		return documentNumber, nil
	}
}

func (t *DocumentNumbersRepositoryImpl) GetAllByUserID(userId string) ([]model.DocumentNumbers, *helper.ErrorModel) {

	helper.PrintValue("Get All By User ID", userId)
	var documentNumbers []model.DocumentNumbers
	result := t.Db.Preload("NumberingFormat").Preload("NumberingFormat.Group").Where("deleted_at IS NULL").Where("state != 0").Where("user_id = ?", userId).Find(&documentNumbers)
	if result.Error != nil {
		msg := "Failed to get all numbering formats"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentNumbers, nil
}

func (t *DocumentNumbersRepositoryImpl) Delete(id string) *helper.ErrorModel {
	documentNumber := &model.DocumentNumbers{}
	documentNumbersId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	deletedRecord := t.Db.Where("id = ?", documentNumbersId).
		First(documentNumber)

	if deletedRecord.Error != nil {
		if errors.Is(deletedRecord.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		msg := "Document Number not found"
		return helper.ErrorCatcher(deletedRecord.Error, 404, &msg)
	}

	documentNumber.State = int(enums.Cancelled)
	result := t.Db.Save(&documentNumber)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		msg := "Document Number not found"
		return helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return nil
}
