package mocks

import "errors"

type DbOperationError int

const (
	DbOpSuccess DbOperationError = iota
	DbOpInternalError
	DbOpDuplicateEmail
	DbOpDuplicatePhone
	DbOpDuplicateFile
	DbOpNotFound
	DbOpInvalidCredentials
)

type GrpcOperationError int

const (
	GrpcOpSuccess GrpcOperationError = iota
	GrpcOpMissingMetadata
	GrpcOpMissingUserID
	GrpcOpInternalError
	GrpcOpInvalidArgument
	GrpcOpFileNameAlreadyExists
	GrpcOpFilePathAlreadyExists
	GrpcOpNotFound
	GrpcOpRollbackFailure
)

var (
	ErrMissingMetadata       = errors.New("missing metadata")
	ErrMissingUserID         = errors.New("missing user id in metadata")
	ErrFileIDRequired        = errors.New("file ID is required")
	ErrFailedToUploadFile    = errors.New("failed to upload file")
	ErrFailedToDownloadFile  = errors.New("failed to download file")
	ErrFailedToRenameFile    = errors.New("failed to rename file")
	ErrFailedToDeleteFile    = errors.New("failed to delete file")
	ErrFileNameAlreadyExists = errors.New("file name already exists")
	ErrFilePathAlreadyExists = errors.New("file path already exists")
	ErrFileNotFound          = errors.New("file not found")
	ErrFailedToGetFiles      = errors.New("failed to get files")
	ErrFailedToGetFileByID   = errors.New("failed to get file by ID")
	ErrFailedToUpdateFile    = errors.New("failed to update file")
)
