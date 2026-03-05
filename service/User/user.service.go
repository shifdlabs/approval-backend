package user

import (
	request "Microservice/data/request/User"
	response "Microservice/data/response/User"
	"Microservice/helper"
)

type UserService interface {
	Create(data request.CreateUserRequest) *helper.ErrorModel
	Get(id string) (*response.UserResponse, *helper.ErrorModel)
	GetAll() ([]response.UserResponse, *helper.ErrorModel)
	GetAllUserExceptCurrent(userId string) ([]response.UserResponse, *helper.ErrorModel)
	Update(request request.UpdateUserRequest) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	MultipleDelete(ids []string) *helper.ErrorModel

	UpdateBiodata(id string, request request.UpdateBiodataRequest) *helper.ErrorModel
	UpdateEmail(id string, request request.UpdateEmailRequest) *helper.ErrorModel
	UpdateRole(request request.UpdateRoleRequest) *helper.ErrorModel
	UpdatePassword(request request.UpdatePasswordRequest) *helper.ErrorModel
	UpdateAccess(request request.UpdateAccessRequest) *helper.ErrorModel
}
