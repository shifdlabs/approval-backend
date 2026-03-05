package documentattachment

import (
	"Microservice/helper"
	"Microservice/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentAttachmentRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentAttachmentRepositoryImpl(Db *gorm.DB) DocumentAttachmentRepository {
	return &DocumentAttachmentRepositoryImpl{Db: Db}
}

func (t *DocumentAttachmentRepositoryImpl) Create(db *gorm.DB, document model.DocumentAttachment) *helper.ErrorModel {
	result := db.Create(&document)

	if result.Error != nil {
		msg := "Failed to create document attachment"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentAttachmentRepositoryImpl) Get(id string) (*model.DocumentAttachment, *helper.ErrorModel) {
	var documentAttachment model.DocumentAttachment
	documentAttachmentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.First(&documentAttachment, "id = ?", documentAttachmentId)

	if strings.Contains(result.Error.Error(), "record not found") {
		return nil, nil
	}

	if result.Error != nil {
		msg := "Get Document attachment failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentAttachment, nil
}

func (t *DocumentAttachmentRepositoryImpl) GetAll() ([]model.DocumentAttachment, *helper.ErrorModel) {
	var documentAttachments []model.DocumentAttachment
	result := t.Db.Find(&documentAttachments)
	if result.Error != nil {
		msg := "Failed to get all documents attachments"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentAttachments, nil
}

func (t *DocumentAttachmentRepositoryImpl) Delete(id string) *helper.ErrorModel {
	documentAttachmentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Unscoped().Delete(&model.DocumentAttachment{}, documentAttachmentId)
	if result.Error != nil {
		msg := "Failed to delete document attachments"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}
