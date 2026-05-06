package user

import (
	"Microservice/helper"
	"Microservice/model"

	request "Microservice/data/request/User"
	response "Microservice/data/response/User"

	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	positionRepository "Microservice/repository/Position"
	repository "Microservice/repository/User"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository               repository.UserRepository
	PositionRepository           positionRepository.PositionRepository
	FailedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository
	Validate                     *validator.Validate
}

func NewUserServiceImpl(
	userRepository repository.UserRepository,
	positionRepository positionRepository.PositionRepository,
	failedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository,
	validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository:               userRepository,
		PositionRepository:           positionRepository,
		FailedLoginAttemptRepository: failedLoginAttemptRepository,
		Validate:                     validate,
	}
}

func (t UserServiceImpl) Create(request request.CreateUserRequest) *helper.ErrorModel {
	var position *model.Position

	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	if request.PositionID != "" {
		result, errGetPosition := t.PositionRepository.Get(request.PositionID)
		if errGetPosition != nil {
			return errGetPosition
		} else {
			position = result
		}
	}

	hashedPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if errBcrypt != nil {
		msg := "Failed to encrypt password"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	newUser := model.User{
		Position:   position,
		EmployeeID: request.EmployeeID,
		Email:      request.Email,
		Password:   string(hashedPassword),
		Role:       request.Role,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Access:     request.Access,
		Phone:      request.Phone,
	}

	errCreateUser := t.UserRepository.Create(newUser)
	if errCreateUser != nil {
		return errCreateUser
	}

	return nil
}

func (t UserServiceImpl) Get(id string) (*response.UserResponse, *helper.ErrorModel) {
	result, errGetUser := t.UserRepository.Get(id, true)
	if errGetUser != nil {
		return nil, errGetUser
	}

	response := response.UserResponse{
		ID:        result.ID,
		Email:     result.Email,
		Role:      result.Role,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Position:  result.Position,
		Access:    result.Access,
		Phone:     result.Phone,
		CreatedAt: *result.CreatedAt,
		UpdatedAt: *result.UpdatedAt,
	}

	return &response, nil
}

func (t UserServiceImpl) GetAll() ([]response.UserResponse, *helper.ErrorModel) {
	response, errGetUsers := t.UserRepository.GetAll()

	if errGetUsers != nil {
		return nil, errGetUsers
	} else {
		return t.mapUsertoUserResponse(response), nil
	}
}

func (t UserServiceImpl) GetAllUserExceptCurrent(userId string) ([]response.UserResponse, *helper.ErrorModel) {
	response, errGetUsers := t.UserRepository.GetAllUserExceptCurrent(userId)

	if errGetUsers != nil {
		return nil, errGetUsers
	} else {
		return t.mapUsertoUserResponse(response), nil
	}
}

func (t UserServiceImpl) Update(request request.UpdateUserRequest) *helper.ErrorModel {
	var position *model.Position

	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	if request.PositionID != "" {
		result, errGetPosition := t.PositionRepository.Get(request.PositionID)
		if errGetPosition != nil {
			return errGetPosition
		} else {
			position = result
		}
	}

	result, errGetUser := t.UserRepository.Get(request.ID, false)
	if errGetUser != nil {
		return errGetUser
	}

	result.FirstName = request.FirstName
	result.LastName = request.LastName
	result.Email = request.Email
	result.Role = request.Role
	result.Position = position
	result.EmployeeID = request.EmployeeID
	result.Access = request.Access
	result.Phone = request.Phone

	// If Access is being enabled, also unlock the account and clear failed attempts
	if request.Access == true && result.IsLocked {
		result.IsLocked = false
		result.LockTimestamp = nil

		// Clear all failed login attempts for this user
		errDeleteAttempts := t.FailedLoginAttemptRepository.DeleteByUserId(result.ID.String())
		if errDeleteAttempts != nil {
			// Log error but continue with update
			helper.GetFileAndLine(errDeleteAttempts)
		}
	}

	errUpdate := t.UserRepository.Update(*result)

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) Delete(id string) *helper.ErrorModel {
	errResponse := t.UserRepository.Delete(id)

	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (t UserServiceImpl) MultipleDelete(ids []string) *helper.ErrorModel {
	errResponse := t.UserRepository.MultipleDelete(ids)

	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (t UserServiceImpl) UpdateEmail(id string, request request.UpdateEmailRequest) *helper.ErrorModel {
	user, errGet := t.UserRepository.Get(id, false)
	if errGet != nil {
		return errGet
	}

	user.Email = request.NewEmail

	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) UpdateBiodata(id string, request request.UpdateBiodataRequest) *helper.ErrorModel {
	var position *model.Position

	user, errGet := t.UserRepository.Get(id, false)
	if errGet != nil {
		return errGet
	}

	if request.PositionID != "" {
		result, errGetPosition := t.PositionRepository.Get(request.PositionID)
		if errGetPosition != nil {
			return errGetPosition
		} else {
			position = result
		}
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Phone = request.Phone
	user.Position = position

	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) UpdateRole(request request.UpdateRoleRequest) *helper.ErrorModel {
	user, errGet := t.UserRepository.Get(request.ID, false)
	if errGet != nil {
		return errGet
	}

	user.Role = request.Role

	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) UpdatePassword(request request.UpdatePasswordRequest) *helper.ErrorModel {
	user, errGet := t.UserRepository.Get(request.ID, false)
	if errGet != nil {
		return errGet
	}

	if request.CurrentPassword != "" {
		errVerify := helper.VerifyPassword(user.Password, request.CurrentPassword)
		if errVerify != nil {
			msg := "incorrect Current Password"
			return helper.ErrorCatcher(errVerify, 404, &msg)
		}
	}

	hashedPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if errBcrypt != nil {
		msg := "Failed to hash password"
		return helper.ErrorCatcher(errBcrypt, 500, &msg)
	}

	user.Password = string(hashedPassword)

	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) UpdateAccess(request request.UpdateAccessRequest) *helper.ErrorModel {
	user, errGet := t.UserRepository.Get(request.ID, false)
	if errGet != nil {
		return errGet
	}

	user.Access = request.Access

	// If Access is being enabled, also unlock the account and clear failed attempts
	if request.Access == true && user.IsLocked {
		user.IsLocked = false
		user.LockTimestamp = nil

		// Clear all failed login attempts for this user
		errDeleteAttempts := t.FailedLoginAttemptRepository.DeleteByUserId(user.ID.String())
		if errDeleteAttempts != nil {
			// Log error but continue with update
			helper.GetFileAndLine(errDeleteAttempts)
		}
	}

	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t UserServiceImpl) UnlockUser(userId string) *helper.ErrorModel {
	// Get user by ID
	user, errGet := t.UserRepository.Get(userId, false)
	if errGet != nil {
		return errGet
	}

	// Reset lock status and enable access
	user.IsLocked = false
	user.LockTimestamp = nil
	user.Access = true

	// Update user in database
	errUpdate := t.UserRepository.Update(*user)
	if errUpdate != nil {
		return errUpdate
	}

	// Delete all failed login attempts for this user
	errDelete := t.FailedLoginAttemptRepository.DeleteByUserId(userId)
	if errDelete != nil {
		return errDelete
	}

	return nil
}
