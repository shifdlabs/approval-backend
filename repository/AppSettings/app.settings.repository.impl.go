package appsettings

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type AppSettingsRepositoryImpl struct {
	Db *gorm.DB
}

func NewAppSettingsRepositoryImpl(Db *gorm.DB) AppSettingsRepository {
	return &AppSettingsRepositoryImpl{Db: Db}
}

func (t *AppSettingsRepositoryImpl) Create(report model.AppSettings) *helper.ErrorModel {
	result := t.Db.Create(&report)

	if result.Error != nil {
		msg := "Create AppSettings Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *AppSettingsRepositoryImpl) GetAll() ([]model.AppSettings, *helper.ErrorModel) {
	var appSettingss []model.AppSettings

	result := t.Db.Find(&appSettingss)
	if result.Error != nil {
		msg := "Failed to Get All AppSettings Data"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return appSettingss, nil
}

func (t *AppSettingsRepositoryImpl) Update(appSettings []model.AppSettings) *helper.ErrorModel {
	trx := t.Db.Begin()
	trx.Begin()

	for _, value := range appSettings {
		var existing model.AppSettings
		err := t.Db.Where("key = ?", value.Key).First(&existing).Error
		if err != nil {
			if errCreate := t.Db.Create(&value).Error; errCreate != nil {
				msg := "Failed to Get All AppSettings Data"
				return helper.ErrorCatcher(errCreate, 500, &msg)
			}
		} else {
			if errUpdate := t.Db.Model(&existing).Updates(value).Error; err != nil {
				msg := "Failed to Get All AppSettings Data"
				return helper.ErrorCatcher(errUpdate, 500, &msg)
			}
		}
	}

	trx.Commit()
	return nil
}
