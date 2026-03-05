package bookmark

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BookmarkRepositoryImpl struct {
	Db *gorm.DB
}

func NewBookmarkRepositoryImpl(Db *gorm.DB) BookmarkRepository {
	return &BookmarkRepositoryImpl{Db: Db}
}

// AddBookmark menambahkan bookmark baru
func (r *BookmarkRepositoryImpl) AddBookmark(userID, documentID uuid.UUID) *helper.ErrorModel {
	bookmark := model.Bookmark{
		BookmarkID: uuid.NewV4(),
		UserID:     userID,
		DocumentID: documentID,
	}
	result := r.Db.Create(&bookmark)
	if result.Error != nil {
		msg := "Failed to add bookmark"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

// RemoveBookmark menghapus bookmark berdasarkan UserID dan DocumentID
func (r *BookmarkRepositoryImpl) RemoveBookmark(userID, documentID uuid.UUID) *helper.ErrorModel {
	result := r.Db.Where("user_id = ? AND document_id = ?", userID, documentID).Delete(&model.Bookmark{})
	if result.Error != nil {
		msg := "Failed to remove bookmark"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

// IsBookmarked memeriksa apakah dokumen sudah di-bookmark oleh user tertentu
func (r *BookmarkRepositoryImpl) IsBookmarked(userID, documentID uuid.UUID) (bool, *helper.ErrorModel) {
	var count int64
	result := r.Db.Model(&model.Bookmark{}).Where("user_id = ? AND document_id = ?", userID, documentID).Count(&count)
	if result.Error != nil {
		msg := "Failed to check bookmark status"
		return false, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return count > 0, nil
}

func (r *BookmarkRepositoryImpl) GetAllBookmarksWithDocuments(userID uuid.UUID) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document
	// Join tabel bookmarks dengan tabel documents
	result := r.Db.
		Model(&model.Document{}).
		Preload("Author").
		Joins("JOIN bookmarks ON bookmarks.document_id = documents.id").
		Where("bookmarks.user_id = ?", userID).
		Find(&documents)

	if result.Error != nil {
		msg := "Failed to fetch bookmarks with documents"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}
