package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
)

type DocumentService interface {
	Create(request request.CreateDocumentRequest) (*model.Document, *helper.ErrorModel)
	GetDocument(id string) (*response.DocumentResponse, *helper.ErrorModel)
	GetDetailDocument(id string, currentUserId string) (*response.DocumentDetailResponse, *helper.ErrorModel)
	GetDetailForEdit(id string) (*response.EditDocumentResponse, *helper.ErrorModel)
	GetAllDocument() ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllReferences(query string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllAuthorization(userId string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllInbox(userId string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllInProgress(userId string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetDocumentStatistics(userId string) (*response.DocumentStatisticResponse, *helper.ErrorModel)
	GetInProgressOverview(userId string) (*response.DocumentInProgressResponse, *helper.ErrorModel)
	GetInProgressOverviewByDocId(documentId string) (*response.DocumentInProgressResponse, *helper.ErrorModel)
	GetRejectedOverview(userId string) (*response.RejectedOverviewResponse, *helper.ErrorModel)
	GetCompletedOverview(userId string) (*response.CompletedOverviewResponse, *helper.ErrorModel)
	Update(request request.UpdateDocumentRequest) (*model.Document, *helper.ErrorModel)
	Authorize(request request.Authorize, userId string) *helper.ErrorModel
	GetCompleteByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetDraftByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetRejectedByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllAuthorDocuments(authorID string) ([]response.DocumentResponse, *helper.ErrorModel)
}
