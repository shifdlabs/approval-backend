package bookmark

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type BookmarkRepository interface {
	AddBookmark(userID, documentID uuid.UUID) *helper.ErrorModel
	RemoveBookmark(userID, documentID uuid.UUID) *helper.ErrorModel
	IsBookmarked(userID, documentID uuid.UUID) (bool, *helper.ErrorModel)
	GetAllBookmarksWithDocuments(userID uuid.UUID) ([]model.Document, *helper.ErrorModel)
}
