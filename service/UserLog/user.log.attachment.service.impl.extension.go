package userlog

import (
	response "Microservice/data/response/UserLog"
	"Microservice/model"
)

func (t UserLogServiceImpl) mapUserLogToUserLogResponse(documentHistories []model.UserLog) []response.UserLogResponse {
	responseDocuments := make([]response.UserLogResponse, len(documentHistories))
	for i, userLog := range documentHistories {
		responseDocuments[i] = t.convertUserLogToUserLogResponse(userLog)
	}
	return responseDocuments
}

func (t UserLogServiceImpl) convertUserLogToUserLogResponse(userLog model.UserLog) response.UserLogResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.UserLogResponse{
		Id:      userLog.ID,
		UserID:  &userLog.UserID,
		Action:  userLog.Action,
		Module:  userLog.Module,
		Log:     userLog.Log,
		LogDate: *userLog.LogDate,
	}

	return responseDocument
}
