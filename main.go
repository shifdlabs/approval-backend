package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"Microservice/config"
	"Microservice/controller"
	"Microservice/model"
	"Microservice/router"

	appSettingsRepository "Microservice/repository/AppSettings"
	bookmarkRepository "Microservice/repository/Bookmark"
	carbonCopiesRepository "Microservice/repository/CarbonCopy"
	documentRepository "Microservice/repository/Document"
	documentAttachmentRepository "Microservice/repository/DocumentAttachment"
	documentHistoryRepository "Microservice/repository/DocumentHistory"
	documentNumbersRepository "Microservice/repository/DocumentNumbers"
	documentReferenceRepository "Microservice/repository/DocumentReference"
	documentSequenceRepository "Microservice/repository/DocumentSequence"
	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	numberingFormatRepository "Microservice/repository/NumberingFormat"
	numberingGroupRepository "Microservice/repository/NumberingGroup"
	positionRepository "Microservice/repository/Position"
	recipientRepository "Microservice/repository/Recipient"
	signatureRepository "Microservice/repository/Signature"
	userRepository "Microservice/repository/User"
	userLogRepository "Microservice/repository/UserLog"

	appSettingService "Microservice/service/AppSettings"
	authService "Microservice/service/Authentication"
	bookmarkService "Microservice/service/Bookmark"
	documentService "Microservice/service/Document"
	documentAttachmentService "Microservice/service/DocumentAttachment"
	documentHistoryService "Microservice/service/DocumentHistory"
	documentNumbersService "Microservice/service/DocumentNumbers"
	documentSequenceService "Microservice/service/DocumentSequence"
	numberingFormatService "Microservice/service/NumberingFormat"
	numberingGroupService "Microservice/service/NumberingGroup"
	positionService "Microservice/service/Position"
	recipientService "Microservice/service/Recipient"
	signatureService "Microservice/service/Signature"
	tokenService "Microservice/service/Token"
	userService "Microservice/service/User"
	userLogService "Microservice/service/UserLog"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func main() {
	envConf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	// Redis
	config.ConnectRedis(&envConf)

	// Database
	db := config.DatabaseConnection(&envConf)
	validate := validator.New()

	println("Message: Migrating Table... ")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Document{})
	db.AutoMigrate(&model.Position{})
	db.AutoMigrate(&model.DocumentHistory{})
	db.AutoMigrate(&model.DocumentSequence{})
	db.AutoMigrate(&model.DocumentAttachment{})
	db.AutoMigrate(&model.Recipient{})
	db.AutoMigrate(&model.AppSettings{})
	db.AutoMigrate(&model.UserLog{})
	db.AutoMigrate(&model.CarbonCopy{})
	db.AutoMigrate(&model.NumberingGroup{})
	db.AutoMigrate(&model.NumberingFormat{})
	db.AutoMigrate(&model.DocumentNumbers{})
	db.AutoMigrate(&model.DocumentReference{})
	db.AutoMigrate(&model.Signature{})
	db.AutoMigrate(&model.NumberingGroup{}, &model.NumberingFormat{})
	db.AutoMigrate(&model.NumberingFormat{}, &model.DocumentNumbers{})
	db.AutoMigrate(&model.User{}, &model.Position{})
	db.AutoMigrate(&model.Document{}, &model.DocumentHistory{})
	db.AutoMigrate(&model.Document{}, &model.DocumentSequence{})
	db.AutoMigrate(&model.Document{}, &model.DocumentAttachment{})
	db.AutoMigrate(&model.Document{}, &model.Recipient{})
	db.AutoMigrate(&model.Document{}, &model.CarbonCopy{})
	db.AutoMigrate(&model.Document{}, &model.Bookmark{})
	db.AutoMigrate(&model.Document{}, &model.DocumentNumbers{})
	db.AutoMigrate(&model.User{}, &model.Signature{})

	// Repositories
	userRepository := userRepository.NewUserRepositoryImpl(db)
	documentRepository := documentRepository.NewDocumentRepositoryImpl(db)
	documentHistoryRepository := documentHistoryRepository.NewDocumentHistoryRepositoryImpl(db)
	documentSequenceRepository := documentSequenceRepository.NewDocumentSequenceRepositoryImpl(db)
	documentAttachmentRepository := documentAttachmentRepository.NewDocumentAttachmentRepositoryImpl(db)
	positionRepositoy := positionRepository.NewPositionRepositoryImpl(db)
	userLogRepository := userLogRepository.NewUserLogRepositoryImpl(db)
	appSettingsRepository := appSettingsRepository.NewAppSettingsRepositoryImpl(db)
	recipientRepository := recipientRepository.NewRecipientRepositoryImpl(db)
	carbonCopiesRepository := carbonCopiesRepository.NewCarbonCopyRepositoryImpl(db)
	bookmarkRepository := bookmarkRepository.NewBookmarkRepositoryImpl(db)
	numberingGroupRepository := numberingGroupRepository.NewNumberingGroupRepositoryImpl(db)
	numberingFormatRepository := numberingFormatRepository.NewNumberingFormatRepositoryImpl(db)
	documentNumbersRepository := documentNumbersRepository.NewDocumentNumbersRepositoryImpl(db)
	documentReferenceRepository := documentReferenceRepository.NewDocumentReferenceRepositoryImpl(db)
	signatureRepository := signatureRepository.NewSignatureRepositoryImpl(db)
	failedLoginAttemptRepository := failedLoginAttemptRepository.NewFailedLoginAttemptRepositoryImpl(db)

	// Servic
	tokenService := tokenService.NewTokenServiceImpl(userRepository)
	authService := authService.NewAuthServiceImpl(userRepository, failedLoginAttemptRepository, validate)
	userService := userService.NewUserServiceImpl(userRepository, positionRepositoy, failedLoginAttemptRepository, validate)
	userLogService := userLogService.NewUserLogServiceImpl(userLogRepository, validate)
	documentSequenceService := documentSequenceService.NewDocumentSequenceServiceImpl(documentSequenceRepository, validate)
	documentService := documentService.NewDocumentServiceImpl(documentRepository, userRepository, documentSequenceRepository, documentAttachmentRepository, documentHistoryRepository, recipientRepository, carbonCopiesRepository, userLogRepository, documentNumbersRepository, documentReferenceRepository, signatureRepository, db, validate)
	documentHistoryService := documentHistoryService.NewDocumentHistoryServiceImpl(documentHistoryRepository, validate)
	documentAttachmentService := documentAttachmentService.NewDocumentAttachmentServiceImpl(documentAttachmentRepository, validate)
	positionService := positionService.NewPositionServiceImpl(positionRepositoy, validate)
	appSettingsService := appSettingService.NewAppSettingsServiceImpl(appSettingsRepository, validate)
	recipientService := recipientService.NewRecipientServiceImpl(recipientRepository, documentRepository, db, validate)
	bookmarkService := bookmarkService.NewBookmarkServiceImpl(bookmarkRepository, validate)
	numberingGroupService := numberingGroupService.NewNumberingGroupServiceImpl(numberingGroupRepository, validate)
	numberingFormatService := numberingFormatService.NewNumberingFormatServiceImpl(numberingFormatRepository, numberingGroupRepository, validate)
	documentNumbersService := documentNumbersService.NewDocumentNumbersServiceImpl(documentNumbersRepository, numberingFormatRepository, validate)
	signatureService := signatureService.NewSignatureServiceImpl(signatureRepository, validate)

	// Controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService, userService)
	tokenController := controller.NewTokenController(tokenService)
	documentController := controller.NewDocumentController(documentService, documentNumbersService, userLogService)
	documentHistoryController := controller.NewDocumentHistoryController(documentHistoryService)
	documentSequenceController := controller.NewDocumentSequenceController(documentSequenceService)
	documentAttachmentController := controller.NewDocumentAttachmentController(documentAttachmentService, userLogService)
	positionController := controller.NewPositionController(positionService, userLogService)
	userLogController := controller.NewUserLogController(userLogService)
	appSettingsController := controller.NewAppSettingsController(appSettingsService)
	recipientController := controller.NewRecipientController(recipientService)
	bookmarkController := controller.NewBookmarkController(bookmarkService)
	numberingGroupController := controller.NewNumberingGroupController(numberingGroupService, userLogService)
	numberingFormatController := controller.NewNumberingFormatController(numberingFormatService, userLogService)
	documentNumbersController := controller.NewDocumentNumbersController(documentNumbersService, userLogService)
	signatureController := controller.NewSignatureController(signatureService)

	// Initialize Router
	routes := router.NewRouter(
		db,
		userController,
		authController,
		tokenController,
		documentController,
		documentHistoryController,
		documentAttachmentController,
		documentSequenceController,
		positionController,
		userLogController,
		appSettingsController,
		recipientController,
		bookmarkController,
		numberingGroupController,
		numberingFormatController,
		documentNumbersController,
		signatureController,
	)

	// Intialize Server
	server := &http.Server{
		Addr:           ":8081",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	errProvideServer := server.ListenAndServe().Error()

	fmt.Println(errProvideServer)

	println("Message: Server Successfully Running...")
}

/*
 Redis should store the User ID, with Key Token.ID or Token
 So when log out, we can extracting the Token.ID or Token as a Key, and get the User ID value for Query processing
 Task:
 - Change Key to Token ID / Token Instead of UserID, in Login and RefreshToken
 -
 Open code for Extracting Token ID value from Redis to get User ID in Logout Function
*/

/*
ChangeLog:
- All Redish process in AUTH flow was removed, because we no longer need it to store an identifier
- Add refresh token expired handler, so when refresh token is expired, logout the user from system
- id in payload access token & refresh token is user id
*/
