package user

import (
	"Microservice/helper"
	"Microservice/model"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	request "Microservice/data/request/User"
	response "Microservice/data/response/User"

	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	positionRepository "Microservice/repository/Position"
	repository "Microservice/repository/User"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
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
		return helper.ErrorCatcher(errStructure, 400, &msg)
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
		return helper.ErrorCatcher(errBcrypt, 500, &msg)
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
		return helper.ErrorCatcher(errStructure, 400, &msg)
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

func (t UserServiceImpl) PreviewImport(fileHeader *multipart.FileHeader, columnMappingJSON string) (*response.PreviewImportResponse, *helper.ErrorModel) {
	// Parse column mapping
	var columnMapping map[string]string
	if err := json.Unmarshal([]byte(columnMappingJSON), &columnMapping); err != nil {
		msg := "Invalid column mapping format"
		return nil, &helper.ErrorModel{Code: 400, Message: msg}
	}

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		msg := "Failed to open file"
		return nil, &helper.ErrorModel{Code: 500, Message: msg}
	}
	defer file.Close()

	// Read Excel file
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		msg := "Failed to read Excel file"
		return nil, &helper.ErrorModel{Code: 400, Message: msg}
	}
	defer excelFile.Close()

	// Get first sheet
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		msg := "Excel file is empty or invalid"
		return nil, &helper.ErrorModel{Code: 400, Message: msg}
	}

	// Get header row and create column index map
	headerRow := rows[0]
	columnIndexes := make(map[string]int)
	for i, header := range headerRow {
		// Store with lowercase and trimmed for case-insensitive matching
		normalizedHeader := strings.ToLower(strings.TrimSpace(header))
		columnIndexes[normalizedHeader] = i
	}

	// Parse data rows
	var users []request.ImportedUserData
	var errors []response.ImportError
	var warnings []response.ImportWarning

	for rowIndex, row := range rows[1:] { // Skip header row
		actualRow := rowIndex + 2 // Excel row number (1-indexed + header)

		user := request.ImportedUserData{}
		hasError := false

		// Map columns based on user-provided mapping
		if colName, ok := columnMapping["employeeID"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.EmployeeID = strings.TrimSpace(row[colIndex])
			}
		}

		if colName, ok := columnMapping["email"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.Email = strings.TrimSpace(row[colIndex])
				if user.Email == "" {
					errors = append(errors, response.ImportError{
						Row:     actualRow,
						Field:   "email",
						Message: "Email is required",
					})
					hasError = true
				}
			} else {
				errors = append(errors, response.ImportError{
					Row:     actualRow,
					Field:   "email",
					Message: "Email column not found",
				})
				hasError = true
			}
		}

		if colName, ok := columnMapping["firstName"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.FirstName = strings.TrimSpace(row[colIndex])
			}
		}

		if colName, ok := columnMapping["lastName"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.LastName = strings.TrimSpace(row[colIndex])
			}
		}

		if colName, ok := columnMapping["phone"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.Phone = strings.TrimSpace(row[colIndex])
			}
		}

		if colName, ok := columnMapping["role"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				roleStr := strings.TrimSpace(row[colIndex])
				if roleInt, err := strconv.Atoi(roleStr); err == nil {
					user.Role = roleInt
				} else {
					warnings = append(warnings, response.ImportWarning{
						Row:     actualRow,
						Field:   "role",
						Message: "Invalid role value",
					})
				}
			}
		}
		// Note: No default role in preview - will be set during bulk insert

		if colName, ok := columnMapping["positionID"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.PositionID = strings.TrimSpace(row[colIndex])
			}
		}

		// Position lookup by name
		if colName, ok := columnMapping["position"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				positionName := strings.TrimSpace(row[colIndex])
				if positionName != "" {
					// Lookup position by name (case-insensitive)
					foundPosition, errPosition := t.PositionRepository.FindByName(positionName)
					if errPosition != nil {
						// Database error - show warning
						warnings = append(warnings, response.ImportWarning{
							Row:     actualRow,
							Field:   "position",
							Message: fmt.Sprintf("Error looking up position '%s'", positionName),
						})
					} else if foundPosition != nil {
						// Position found, use its ID
						user.PositionID = foundPosition.ID.String()
					} else {
						// Position not found, silently default to "Staff"
						staffPosition, errStaff := t.PositionRepository.FindByName("Staff")
						if errStaff == nil && staffPosition != nil {
							user.PositionID = staffPosition.ID.String()
							// No warning - silent fallback to Staff
						} else {
							// Only warn if Staff position doesn't exist (critical case)
							warnings = append(warnings, response.ImportWarning{
								Row:     actualRow,
								Field:   "position",
								Message: "No 'Staff' position found in database for default",
							})
						}
					}
				}
			}
		}

		// Password defaults will be handled during bulk insert
		if colName, ok := columnMapping["password"]; ok && colName != "" {
			normalizedColName := strings.ToLower(strings.TrimSpace(colName))
			if colIndex, exists := columnIndexes[normalizedColName]; exists && colIndex < len(row) {
				user.Password = strings.TrimSpace(row[colIndex])
			}
		}

		if !hasError {
			users = append(users, user)
		}
	}

	return &response.PreviewImportResponse{
		Users:    users,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

func (t UserServiceImpl) BulkImport(importRequest request.BulkImportUsersRequest) (*response.BulkImportResponse, *helper.ErrorModel) {
	errStructure := t.Validate.Struct(importRequest)
	if errStructure != nil {
		msg := "Structure Error"
		return nil, helper.ErrorCatcher(errStructure, 400, &msg)
	}

	var successCount, failedCount int
	var errors []response.ImportError

	fmt.Printf("Starting bulk import for %d users\n", len(importRequest.Users))

	for i, userData := range importRequest.Users {
		rowNum := i + 2 // Excel row number (1-indexed + header)

		fmt.Printf("Processing row %d: %s\n", rowNum, userData.Email)

		// Validate email
		if userData.Email == "" {
			errors = append(errors, response.ImportError{
				Row:     rowNum,
				Field:   "email",
				Message: "Email is required",
			})
			failedCount++
			fmt.Printf("Row %d failed: Email is required\n", rowNum)
			continue
		}

		// Check if email already exists
		existingUser, _ := t.UserRepository.GetByEmail(userData.Email)
		if existingUser != nil {
			errors = append(errors, response.ImportError{
				Row:     rowNum,
				Field:   "email",
				Message: fmt.Sprintf("Email %s already exists", userData.Email),
			})
			failedCount++
			fmt.Printf("Row %d failed: Email %s already exists\n", rowNum, userData.Email)
			continue
		}

		// Get position if PositionID provided
		var position *model.Position
		if userData.PositionID != "" {
			fmt.Printf("Looking up position ID: %s\n", userData.PositionID)
			result, errGetPosition := t.PositionRepository.Get(userData.PositionID)
			if errGetPosition != nil {
				errors = append(errors, response.ImportError{
					Row:     rowNum,
					Field:   "positionID",
					Message: fmt.Sprintf("Position ID %s not found or invalid: %s", userData.PositionID, errGetPosition.Message),
				})
				failedCount++
				fmt.Printf("Row %d failed: Position ID lookup error: %s\n", rowNum, errGetPosition.Message)
				continue
			}
			position = result
			fmt.Printf("Position found: %s\n", position.Name)
		}

		password := userData.Password
		if importRequest.CustomPassword != "" {
			password = importRequest.CustomPassword
		} else if password == "" {
			password = generateRandomPassword()
		}

		role := userData.Role
		if role == 0 {
			role = 1
		}

		hashedPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if errBcrypt != nil {
			errors = append(errors, response.ImportError{
				Row:     rowNum,
				Field:   "password",
				Message: "Failed to encrypt password",
			})
			failedCount++
			continue
		}

		// Create new user
		newUser := model.User{
			Position:   position,
			EmployeeID: userData.EmployeeID,
			Email:      userData.Email,
			Password:   string(hashedPassword),
			Role:       role,
			FirstName:  userData.FirstName,
			LastName:   userData.LastName,
			Access:     true, // Default to true
			Phone:      userData.Phone,
		}

		errCreateUser := t.UserRepository.Create(newUser)
		if errCreateUser != nil {
			errors = append(errors, response.ImportError{
				Row:     rowNum,
				Field:   "general",
				Message: errCreateUser.Message,
			})
			failedCount++
			continue
		}

		successCount++
	}

	return &response.BulkImportResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Errors:       errors,
	}, nil
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

// generateRandomPassword generates a cryptographically secure random password.
func generateRandomPassword() string {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("tmp%d", 0)
	}
	return base64.URLEncoding.EncodeToString(b)[:24]
}
