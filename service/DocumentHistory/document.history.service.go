package documentHistory

import (
	response "Microservice/data/response/DocumentHistory"
	"Microservice/helper"
)

type DocumentHistoryService interface {
	Get(id string) (*response.DocumentHistoryResponse, *helper.ErrorModel)
	GetAll() ([]response.DocumentHistoryResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
	FetchHistoriesByUserID(userID string) ([]response.DocumentHistoryResponse, *helper.ErrorModel)
}
