package handlers

import (
	"github.com/gin-gonic/gin"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (
	// functions name for upload user files
	fnUploadFile = "UploadFile"

	// functions name for download user files
	fnDownloadFile = "DownloadFile"

	// functions name for list user files
	fnListFiles = "ListFiles"

	// functions name for rename user files
	fnRenameFile = "RenameFile"

	// functions name for delete user files
	fnDeleteFile = "DeleteFile"
)

type UserFilesHandler struct {
	storage   storage.UserFilesService
	MMSClient MMSpb.FilesClient
}

func NewUserFilesHandler(storage storage.UserFilesService, MMSClient MMSpb.FilesClient) *UserFilesHandler {
	return &UserFilesHandler{storage: storage, MMSClient: MMSClient}
}

func RegisterRoutes(rg *gin.RouterGroup, h *UserFilesHandler) {
	rg.Use(middleware.AuthMiddleware())
	rg.POST("/upload", h.UploadFileHandler)
	rg.GET("/:fileID/download", h.DownloadFileHandler)
	rg.GET("", h.ListFilesHandler)
	rg.PATCH("/:fileID/rename", h.RenameFileHandler)
	rg.DELETE("/:fileID", h.DeleteFileHandler)
}
