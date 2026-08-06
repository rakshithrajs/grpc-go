package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (
	// FnUploadFile is the log prefix for the upload handler.
	FnUploadFile = "UploadFile"

	// FnDownloadFile is the log prefix for the download handler.
	FnDownloadFile = "DownloadFile"

	// FnListFiles is the log prefix for the list files handler.
	FnListFiles = "ListFiles"

	// FnRenameFile is the log prefix for the rename handler.
	FnRenameFile = "RenameFile"

	// FnDeleteFile is the log prefix for the delete handler.
	FnDeleteFile = "DeleteFile"

	multipartFileField = "file"
)

// UserFilesHandler exposes HTTP handlers for user file operations.
type UserFilesHandler struct {
	client  *grpcClient.MMSClient
	storage storage.UserFilesService
}

// NewUserFilesHandler creates a new UserFilesHandler.
func NewUserFilesHandler(client *grpcClient.MMSClient, storage storage.UserFilesService) *UserFilesHandler {
	return &UserFilesHandler{client: client, storage: storage}
}

// RegisterRoutes wires the file operation routes to the provided router group.
func RegisterRoutes(rg *gin.RouterGroup, h *UserFilesHandler, tmsClient *grpcClient.TMSClient) {
	filesRouterGroup := rg.Group("/files")
	filesRouterGroup.Use(middleware.AuthMiddleware(tmsClient))
	filesRouterGroup.POST("/upload", h.UploadFileHandler)
	filesRouterGroup.GET("/:fileid/download", h.DownloadFileHandler)
	filesRouterGroup.GET(config.NullString, h.ListFilesHandler)
	filesRouterGroup.PATCH("/:fileid/rename", h.RenameFileHandler)
	filesRouterGroup.DELETE("/:fileid", h.DeleteFileHandler)
}
