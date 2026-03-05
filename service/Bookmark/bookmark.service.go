package bookmark

import (
	request "Microservice/data/request/Bookmark"
	"Microservice/helper"
	"Microservice/model"
)

type BookmarkService interface {
	AddBookmark(request request.BookmarkRequest) *helper.ErrorModel
	RemoveBookmark(request request.BookmarkRequest) *helper.ErrorModel
	IsBookmarked(request request.BookmarkRequest) (bool, *helper.ErrorModel)
	GetAllBookmarksWithDocuments(userID string) ([]model.Document, *helper.ErrorModel)
}
