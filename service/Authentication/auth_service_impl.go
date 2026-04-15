package authentication

import (
	model "Microservice/data/model/Authentication"
	authentication "Microservice/data/request/Authentication"
	"Microservice/helper"
	dbModel "Microservice/model"
	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	userRepository "Microservice/repository/User"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthServiceImpl struct {
	UserRepository               userRepository.UserRepository
	FailedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository
	Validate                     *validator.Validate
}

func NewAuthServiceImpl(userRepository userRepository.UserRepository, failedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository:               userRepository,
		FailedLoginAttemptRepository: failedLoginAttemptRepository,
		Validate:                     validate,
	}
}

func (t AuthServiceImpl) Login(payload authentication.LogInRequest) (model.LoginResult, *helper.ErrorModel) {
	// Step 1: Get user by email
	user, err := t.UserRepository.GetByEmail(payload.Email)
	if err != nil {
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, err
	}

	// Step 2: Check Access field first
	if !user.Access {
		msg := "Your account access has been disabled. Please contact administrator."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, &helper.ErrorModel{Code: 403, Message: msg}
	}

	// Step 3: Check if account is locked
	if user.IsLocked {
		msg := "Your account is locked due to multiple failed login attempts. Please contact administrator."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, &helper.ErrorModel{Code: 403, Message: msg}
	}

	// Step 4: Verify password
	errVerifyPassword := helper.VerifyPassword(user.Password, payload.Password)
	if errVerifyPassword != nil {
		// Password is incorrect, record failed attempt
		now := time.Now()
		failedAttempt := dbModel.FailedLoginAttempt{
			UserID:      user.ID,
			AttemptedAt: &now,
		}

		// Save failed attempt to database
		errCreate := t.FailedLoginAttemptRepository.Create(failedAttempt)
		if errCreate != nil {
			// Log error but continue with login flow
			helper.GetFileAndLine(errCreate)
		}

		// Count total failed attempts for this user
		count, errCount := t.FailedLoginAttemptRepository.CountByUserId(user.ID.String())
		if errCount != nil {
			msg := "Incorrect password"
			return model.LoginResult{
				AccessToken:  "",
				RefreshToken: "",
				User:         nil,
			}, helper.ErrorCatcher(errVerifyPassword, 400, &msg)
		}

		// Check if attempts >= 3, lock the account
		if count >= 3 {
			// Lock the user account and disable access
			lockTime := time.Now()
			user.IsLocked = true
			user.LockTimestamp = &lockTime
			user.Access = false

			// Update user in database
			errUpdate := t.UserRepository.Update(*user)
			if errUpdate != nil {
				// Log error but still return locked message
				helper.GetFileAndLine(errUpdate)
			}

			msg := "Too many failed login attempts. Your account has been locked. Please contact administrator."
			return model.LoginResult{
				AccessToken:  "",
				RefreshToken: "",
				User:         nil,
			}, helper.ErrorCatcher(errVerifyPassword, 403, &msg)
		}

		// Return error with remaining attempts
		remainingAttempts := 3 - count
		msg := strconv.Itoa(int(remainingAttempts)) + " attempt(s) remaining before account lock."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, helper.ErrorCatcher(errVerifyPassword, 400, &msg)
	}

	// Step 5: Password is correct - Clear failed attempts and unlock if needed
	errDelete := t.FailedLoginAttemptRepository.DeleteByUserId(user.ID.String())
	if errDelete != nil {
		// Log error but continue with successful login
		helper.GetFileAndLine(errDelete)
	}

	// Reset lock status and access if user was locked
	if user.IsLocked {
		user.IsLocked = false
		user.LockTimestamp = nil
		user.Access = true
		errUpdate := t.UserRepository.Update(*user)
		if errUpdate != nil {
			// Log error but continue with successful login
			helper.GetFileAndLine(errUpdate)
		}
	}

	// Step 6: Generate tokens and return success
	accessToken, _ := helper.GenerateAccessToken(user.ID.String())
	refreshToken, _ := helper.GenerateRefreshToken(user.ID.String())

	return model.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// func (t AuthServiceImpl) CheckRegisteredEmail(payload authentication.VerifyForgetPassword) bool {
// 	user, _ := t.UserRepository.GetByEmail(payload.Email)

// 	if user != nil {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (t AuthServiceImpl) ResetPassword(payload authentication.ResetPassword) *helper.CustomError {
// 	err := t.UserRepository.UpdatePasssword(payload.Email, payload.NewPassword)

// 	if err != nil {
// 		return err
// 	} else {
// 		return nil
// 	}
// }
