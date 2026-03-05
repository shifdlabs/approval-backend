package numberingformat

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type NumberingFormatRepositoryImpl struct {
	Db *gorm.DB
}

func NewNumberingFormatRepositoryImpl(Db *gorm.DB) NumberingFormatRepository {
	return &NumberingFormatRepositoryImpl{Db: Db}
}

func (t *NumberingFormatRepositoryImpl) Create(data model.NumberingFormat) *helper.ErrorModel {
	result := t.Db.Create(&data)

	if result.Error != nil {
		msg := "Failed to create numbering format"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *NumberingFormatRepositoryImpl) Get(id string) (*model.NumberingFormat, *helper.ErrorModel) {
	// Return nil if ID is empty
	if id == "" {
		return nil, nil
	}

	var numberingFormat model.NumberingFormat
	numberingFormatId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Where("id = ?", numberingFormatId).First(&numberingFormat)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "CarbonCopys not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &numberingFormat, nil
}

func (t *NumberingFormatRepositoryImpl) GetAll() ([]model.NumberingFormat, *helper.ErrorModel) {
	var numberingFormats []model.NumberingFormat
	result := t.Db.Preload("Group").Where("deleted_at IS NULL").Find(&numberingFormats)
	if result.Error != nil {
		msg := "Failed to get all numbering formats"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return numberingFormats, nil
}

func (t *NumberingFormatRepositoryImpl) Delete(id string) *helper.ErrorModel {
	numberingFormatId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Delete(&model.NumberingFormat{}, numberingFormatId)
	if result.Error != nil {
		msg := "Failed to delete numbering format"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}
