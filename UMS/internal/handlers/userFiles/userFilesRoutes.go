package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (
	// functions name for upload user files
	FnUploadFile = "UploadFile"

	// functions name for download user files
	FnDownloadFile = "DownloadFile"

	// functions name for list user files
	FnListFiles = "ListFiles"

	// functions name for rename user files
	FnRenameFile = "RenameFile"

	// functions name for delete user files
	FnDeleteFile = "DeleteFile"

	multipartFileField = "file"
)

type UserFilesHandler struct {
	client  *grpc.Client
	storage storage.UserFilesService
}

func NewUserFilesHandler(client *grpc.Client, storage storage.UserFilesService) *UserFilesHandler {
	return &UserFilesHandler{client: client, storage: storage}
}

func RegisterRoutes(rg *gin.RouterGroup, h *UserFilesHandler) {
	filesRouterGroup := rg.Group("/files")
	filesRouterGroup.Use(middleware.AuthMiddleware())
	filesRouterGroup.POST("/upload", h.UploadFileHandler)
	filesRouterGroup.GET("/:fileID/download", h.DownloadFileHandler)
	filesRouterGroup.GET(config.NullString, h.ListFilesHandler)
	filesRouterGroup.PATCH("/:fileID/rename", h.RenameFileHandler)
	filesRouterGroup.DELETE("/:fileID", h.DeleteFileHandler)
}
