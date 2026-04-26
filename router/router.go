package router

import (
	"Microservice/config/middleware"
	"Microservice/controller"
	"Microservice/helper"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:5173" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization, X-Requested-With",
			)
			c.Writer.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func NewRouter(
	Db *gorm.DB,
	userController *controller.UserController,
	authController *controller.AuthController,
	tokenController *controller.TokenController,
	documentController *controller.DocumentController,
	documentHistoryController *controller.DocumentHistoryController,
	documentAttachmentController *controller.DocumentAttachmentController,
	documentSequenceController *controller.DocumentSequenceController,
	positionController *controller.PositionController,
	userLogController *controller.UserLogController,
	appSettingsController *controller.AppSettingsController,
	recipientController *controller.RecipientController,
	bookmarkController *controller.BookmarkController,
	numberingGroupController *controller.NumberingGroupController,
	numberingFormatController *controller.NumberingFormatController,
	documentNumbersController *controller.DocumentNumbersController,
	signatureController *controller.SignatureController,
) *gin.Engine {
	service := gin.Default()
	service.Use(CORS())

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "Router has initialized")
	})

	service.NoRoute(func(c *gin.Context) {
		helper.ResponseError(c, helper.CustomError{
			Code:    404,
			Message: "Not Found.",
		})
	})

	router := service.Group("/api")

	authRouter := router.Group("/auth")
	authRouter.POST("/login", authController.LogIn)
	authRouter.POST("/register", authController.Register)
	// authRouter.POST("/forgetPassword", authController.VerifyForgetPassword)
	// authRouter.PUT("/resetPassword", authController.ResetPassword)

	token := router.Group("/refresh")
	token.POST("", tokenController.RefreshAccessToken)

	authRouter.Use(middleware.DeserializeUser(Db))
	authRouter.GET("/logout", authController.Logout)

	protectedUserRouter := router.Group("/user")
	protectedUserRouter.Use(middleware.AuthMiddleware())
	protectedUserRouter.POST("", userController.Create)
	protectedUserRouter.GET("/profile", userController.Get)
	protectedUserRouter.GET("/:id", userController.GetUserByID)
	protectedUserRouter.GET("", userController.GetAll)
	protectedUserRouter.GET("/except-current", userController.GetAllUserExceptCurrent)
	protectedUserRouter.PUT("", userController.Update)
	protectedUserRouter.DELETE("/:id", userController.Delete)
	protectedUserRouter.DELETE("/deletes", userController.MultipleDelete)
	protectedUserRouter.PUT("/role", userController.UpdateRole)
	protectedUserRouter.PUT("/password", userController.UpdatePassword)
	protectedUserRouter.PUT("/access", userController.UpdateAccess)
	protectedUserRouter.PUT("/biodata", userController.UpdateBiodata)
	protectedUserRouter.PUT("/email", userController.UpdateEmail)

	protectedDocumentRouter := router.Group("/document")
	protectedDocumentRouter.Use(middleware.DeserializeUser(Db))
	protectedDocumentRouter.POST("", documentController.Create)
	protectedDocumentRouter.PUT("", documentController.Update)
	protectedDocumentRouter.GET("", documentController.GetAll)
	protectedDocumentRouter.GET("/references/:q", documentController.GetAllReferences)
	protectedDocumentRouter.GET("/:id", documentController.Get)
	protectedDocumentRouter.GET("/detail/:id", documentController.GetDetailPreview)
	protectedDocumentRouter.GET("/edit/:id", documentController.GetDetailForEdit)
	protectedDocumentRouter.POST("/authorize", documentController.Authorize)
	protectedDocumentRouter.GET("/authorization", documentController.GetAllAuthorization)
	protectedDocumentRouter.GET("/inprogress", documentController.GetAllInProgress)
	protectedDocumentRouter.GET("/inbox", documentController.GetAllInbox)
	protectedDocumentRouter.GET("/rejected", documentController.GetAllRejected)
	protectedDocumentRouter.GET("/dashboard", documentController.GetDashboardSummary)
	protectedDocumentRouter.GET("/dashboard/deadlines", documentController.GetDeadlines)
	protectedDocumentRouter.GET("/dashboard/activities", documentController.GetRecentActivities)
	protectedDocumentRouter.GET("/dashboard/recent", documentController.GetRecentDocuments)

	protectedDocumentHistoryRouter := router.Group("/documenthistory")
	protectedDocumentHistoryRouter.Use(middleware.DeserializeUser(Db))
	protectedDocumentHistoryRouter.GET("", documentHistoryController.GetAll)
	protectedDocumentHistoryRouter.GET("/:id", documentHistoryController.Get)
	protectedDocumentHistoryRouter.GET("/rejected", documentHistoryController.GetRejectedWithDocumentAndUser)

	protectedDocumentAttachmentRouter := router.Group("/documentattachment")
	protectedDocumentAttachmentRouter.Use(middleware.DeserializeUser(Db))
	protectedDocumentAttachmentRouter.GET("", documentAttachmentController.GetAll)
	protectedDocumentAttachmentRouter.GET("/:id", documentAttachmentController.Get)
	protectedDocumentAttachmentRouter.DELETE("", documentAttachmentController.Delete)

	protectedDocumentRouter.GET("/complete", documentController.GetComplete)
	protectedDocumentRouter.GET("/draft", documentController.GetDraft)

	protectedUserLogRouter := router.Group("/userlogs")
	protectedUserLogRouter.Use(middleware.DeserializeUser(Db))
	protectedUserLogRouter.GET("", userLogController.GetAll)

	protectedDocumentSequenceRouter := router.Group("/documentsequence")
	protectedDocumentSequenceRouter.Use(middleware.DeserializeUser(Db))
	//protectedDocumentSequenceRouter.GET("", documentSequenceController.GetAll)
	protectedDocumentSequenceRouter.GET("/:id", documentSequenceController.Get)
	protectedDocumentSequenceRouter.GET("/progress", documentSequenceController.GetProgress)

	protectedAppSettingsRouter := router.Group("/appsettings")
	protectedAppSettingsRouter.Use(middleware.DeserializeUser(Db))
	protectedAppSettingsRouter.GET("", appSettingsController.GetAll)
	protectedAppSettingsRouter.PUT("", appSettingsController.Update)

	protectedPositionRouter := router.Group("/position")
	protectedPositionRouter.Use(middleware.DeserializeUser(Db))
	protectedPositionRouter.GET("", positionController.GetAll)
	protectedPositionRouter.GET("/:id", positionController.Get)
	protectedPositionRouter.PUT("", positionController.Update)
	protectedPositionRouter.POST("", positionController.Create)
	protectedPositionRouter.DELETE("/:id", positionController.Delete)

	protectedBookmarkRouter := router.Group("/bookmark")
	protectedBookmarkRouter.Use(middleware.DeserializeUser(Db))
	protectedBookmarkRouter.POST("/add", bookmarkController.AddBookmarkHandler)
	protectedBookmarkRouter.POST("/remove", bookmarkController.RemoveBookmarkHandler)
	protectedBookmarkRouter.POST("/status", bookmarkController.IsBookmarkedHandler)
	protectedBookmarkRouter.GET("/documents", bookmarkController.GetAllBookmarksWithDocumentsHandler)

	protectedNumberingGroupRouter := router.Group("/numbering/group")
	protectedNumberingGroupRouter.Use(middleware.DeserializeUser(Db))
	protectedNumberingGroupRouter.GET("", numberingGroupController.GetAll)
	protectedNumberingGroupRouter.GET("/:id", numberingGroupController.Get)
	protectedNumberingGroupRouter.POST("", numberingGroupController.Create)
	protectedNumberingGroupRouter.DELETE("/:id", numberingGroupController.Delete)

	protectedNumberingFormatRouter := router.Group("/numbering/format")
	protectedNumberingFormatRouter.Use(middleware.DeserializeUser(Db))
	protectedNumberingFormatRouter.GET("", numberingFormatController.GetAll)
	protectedNumberingFormatRouter.GET("/grouped", numberingFormatController.GetAllWithGrouped)
	protectedNumberingFormatRouter.POST("", numberingFormatController.Create)
	protectedNumberingFormatRouter.DELETE("/:id", numberingFormatController.Delete)

	protectedDocumentNumberRouter := router.Group("/document/number")
	protectedDocumentNumberRouter.Use(middleware.DeserializeUser(Db))
	protectedDocumentNumberRouter.POST("", documentNumbersController.Create)
	protectedDocumentNumberRouter.GET("", documentNumbersController.GetAll)
	protectedDocumentNumberRouter.GET("/user", documentNumbersController.GetAllByUserId)
	protectedDocumentNumberRouter.DELETE("/:id", documentNumbersController.Delete)

	protectedSignatureRouter := router.Group("/signature")
	protectedSignatureRouter.Use(middleware.DeserializeUser(Db))
	protectedSignatureRouter.GET("", signatureController.GetAll)
	protectedSignatureRouter.POST("", signatureController.Create)
	protectedSignatureRouter.PUT("/:userId", signatureController.Update)
	protectedSignatureRouter.DELETE("/:userId", signatureController.Delete)
	protectedSignatureRouter.GET("/:userId", signatureController.GetByUserId)

	return service
}
