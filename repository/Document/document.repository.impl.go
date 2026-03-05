package document

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentRepositoryImpl(Db *gorm.DB) DocumentRepository {
	return &DocumentRepositoryImpl{Db: Db}
}

func (t *DocumentRepositoryImpl) Create(db gorm.DB, document *model.Document) *helper.ErrorModel {
	result := db.Create(document)
	if result.Error != nil {
		msg := "Create Document Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) Get(id string) (*model.Document, *helper.ErrorModel) {
	var report model.Document

	reportId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Preload("Author").Preload("DocumentAttachment").Preload("DocumentSequence").Preload("DocumentHistory").First(&report, "id = ?", reportId)

	if result.Error != nil {
		msg := "Get Document Failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &report, nil
}

func (t *DocumentRepositoryImpl) GetAll() ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.Preload("Author").Preload("DocumentAttachment").Preload("DocumentSequence").Preload("DocumentHistory").Preload("DocumentHistory").Find(&reports)
	if result.Error != nil {
		msg := "Failed to get all documents"
		return reports, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetAllReferences(query string) ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.Where("status = ? AND subject ILIKE ?", 2, "%"+query+"%").Find(&reports)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get references"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetAllAuthorization(id string) ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.
		Preload("Author").
		Preload("DocumentHistory").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("document_sequences.user_id = ?", id).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1").
		Find(&reports)

	if result.Error != nil {
		msg := "Failed to get all documents"
		return reports, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetDocumentStatistics(id string) ([]int, *helper.ErrorModel) {
	var totalAuthorization int64

	var totalInProgressAsApprover int64
	var totalInProgressAsAuthor int64

	var totalRejectedAsApprover int64
	var totalRejectedAsAuthor int64

	var totalCompletedAsApprover int64
	var totalCompletedAsAuthor int64

	countAuthorization := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("document_sequences.user_id = ?", id).
		Where("document_sequences.step = documents.step").
		Count(&totalAuthorization)

	if countAuthorization.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countAuthorization.Error, 500, &msg)
	}

	countInProgressAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 1 OR documents.status = 99").
		Count(&totalInProgressAsApprover)

	if countInProgressAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countInProgressAsApprover.Error, 500, &msg)
	}

	countInProgressAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("status = 1").
		Where("author_id = ?", id).
		Count(&totalInProgressAsAuthor)

	if countInProgressAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countInProgressAsAuthor.Error, 500, &msg)
	}

	countRejectedAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 99").
		Count(&totalRejectedAsApprover)

	if countRejectedAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countRejectedAsApprover.Error, 500, &msg)
	}

	countRejectedAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("status = 99").
		Where("author_id = ?", id).
		Count(&totalRejectedAsAuthor)

	if countRejectedAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countRejectedAsAuthor.Error, 500, &msg)
	}

	countCompletedAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 2 OR documents.status = 3").
		Count(&totalCompletedAsApprover)

	if countCompletedAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countCompletedAsApprover.Error, 500, &msg)
	}

	countCompletedAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("status = 2 OR status = 3").
		Where("author_id = ?", id).
		Count(&totalCompletedAsAuthor)

	if countCompletedAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countCompletedAsAuthor.Error, 500, &msg)
	}

	var result = []int{
		int(totalAuthorization),
		int(totalInProgressAsApprover + totalInProgressAsAuthor),
		int(totalRejectedAsApprover + totalRejectedAsAuthor),
		int(totalCompletedAsApprover + totalCompletedAsAuthor),
	}

	return result, nil
}

func (t *DocumentRepositoryImpl) GetOneLatestInprogress(id string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("status = 99").
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetLastestRejected(id string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("status = 99").
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetLastestCompleted(id string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("status = 2").
		Or("status = 3").
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetAllInbox(id string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	response := t.Db.
		Model(&model.Document{}).
		Preload("Author"). // Preload relasi Author jika diperlukan
		Joins("JOIN recipients ON recipients.document_id = documents.id").
		Where("recipients.user_id = ?", id).
		Where("documents.status = 2").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Kembalikan slice kosong jika tidak ada data
			return []model.Document{}, nil
		}

		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetAllInProgress(userId string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	response := t.Db.
		Model(&model.Document{}).
		Select("DISTINCT documents.*").
		Preload("Author"). // Preload relasi Author jika diperlukan
		Preload("DocumentSequence").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.author_id = ? AND documents.status = 1", userId). // first condition
		Or("document_sequences.user_id = ? AND document_sequences.step > documents.step AND documents.status = 1", userId).
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Kembalikan slice kosong jika tidak ada data
			return []model.Document{}, nil
		}

		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) Update(report model.Document) *helper.ErrorModel {
	err := t.Db.Save(&report).Error

	if err != nil {
		msg := "Failed to update document"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) Delete(id string) *helper.ErrorModel {
	reportId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Unscoped().Delete(&model.Document{}, reportId)

	if result.Error != nil {
		msg := "Failed to delete document"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) GetCompleteByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 2
	result := t.Db.Preload("Author").
		Where("author_id = ? AND status = ?", authorID, 2).
		Order("updated_at DESC").
		Find(&documents)
	if result.Error != nil {
		msg := "Failed to fetch documents for the author"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetDraftByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	result := t.Db.Preload("Author").
		Where("author_id = ? AND status = ?", authorID, 0).
		Order("updated_at DESC").
		Find(&documents)
	if result.Error != nil {
		msg := "Failed to fetch draft documents for the author"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetRejectedByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	response := t.Db.Preload("Author").
		Where("author_id = ? AND status = ?", authorID, 99).
		Order("updated_at DESC").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetAllAuthorDocuments(authorID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	response := t.Db.Preload("Author").
		Where("author_id = ?", authorID).
		Order("updated_at DESC").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}
