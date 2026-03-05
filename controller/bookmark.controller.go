package controller

import (
	"Microservice/helper"
	"Microservice/utils"

	request "Microservice/data/request/Bookmark"
	service "Microservice/service/Bookmark"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	bookmarkService service.BookmarkService
}

func NewBookmarkController(service service.BookmarkService) *BookmarkController {
	return &BookmarkController{bookmarkService: service}
}

// AddBookmarkHandler menangani penambahan bookmark
func (controller *BookmarkController) AddBookmarkHandler(ctx *gin.Context) {
	var payload request.BookmarkRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.bookmarkService.AddBookmark(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

// RemoveBookmarkHandler menangani penghapusan bookmark
func (controller *BookmarkController) RemoveBookmarkHandler(ctx *gin.Context) {
	var payload request.BookmarkRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	err := controller.bookmarkService.RemoveBookmark(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, nil)
	}
}

// IsBookmarkedHandler memeriksa apakah dokumen sudah di-bookmark
func (controller *BookmarkController) IsBookmarkedHandler(ctx *gin.Context) {
	var payload request.BookmarkRequest
	errBindJSON := ctx.ShouldBindJSON(&payload)

	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	isBookmarked, err := controller.bookmarkService.IsBookmarked(payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, gin.H{
			"isBookmarked": isBookmarked,
		})
	}
}

func (controller *BookmarkController) GetAllBookmarksWithDocumentsHandler(ctx *gin.Context) {
	// Ambil userId menggunakan helper
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		msg := "Invalid Request Structure."
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	// Panggil service untuk mendapatkan dokumen berdasarkan userId
	documents, err := controller.bookmarkService.GetAllBookmarksWithDocuments(*id)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
	} else {
		utils.SuccessResponse(ctx, documents)
	}
}
