package user

import (
	"Microservice/helper"
	"Microservice/model"
)

type UserRepository interface {
	Create(user model.User) *helper.ErrorModel
	Get(id string, hidePassword bool) (*model.User, *helper.ErrorModel)
	GetAll() ([]model.User, *helper.ErrorModel)
	GetAllUserExceptCurrent(userId string) ([]model.User, *helper.ErrorModel)
	GetByEmail(email string) (*model.User, *helper.ErrorModel)
	Update(user model.User) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	MultipleDelete(ids []string) *helper.ErrorModel
}
