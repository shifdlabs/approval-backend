package document

import (
	dashboard "Microservice/data/model/Dashboard"
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(db gorm.DB, report *model.Document) *helper.ErrorModel
	Get(id string) (*model.Document, *helper.ErrorModel)
	GetAll() ([]model.Document, *helper.ErrorModel)
	GetAllReferences(query string) ([]model.Document, *helper.ErrorModel)
	GetAllAuthorization(id string) ([]model.Document, *helper.ErrorModel)
	GetAllInbox(id string) ([]model.Document, *helper.ErrorModel)
	GetAllInProgress(userId string) ([]model.Document, *helper.ErrorModel)
	GetDocumentStatistics(id string) ([]int, *helper.ErrorModel)
	GetOneLatestInprogress(id string) (*model.Document, *helper.ErrorModel)
	GetLastestRejected(id string) (*model.Document, *helper.ErrorModel)
	GetLastestCompleted(id string) (*model.Document, *helper.ErrorModel)
	Update(report model.Document) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	GetCompleteByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel)
	GetDraftByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel)
	GetRejectedByAuthorID(authorID string) ([]model.Document, *helper.ErrorModel)
	GetAllAuthorDocuments(authorID string) ([]model.Document, *helper.ErrorModel)
	GetDashboardSummary(userId string, period string) (*dashboard.DashboardSummaryRaw, *helper.ErrorModel)
	GetDeadlines(userId string) ([]model.Document, *helper.ErrorModel)
	GetRecentActivities(userId string) ([]model.DocumentHistory, *helper.ErrorModel)
	GetRecentDocuments(userId string, docType int) ([]model.Document, *helper.ErrorModel)
}
