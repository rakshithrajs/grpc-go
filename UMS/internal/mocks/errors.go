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
	ErrFileNameAlreadyExists = errors.New("file name already exists")
	ErrFailedToGetFiles      = errors.New("failed to get files")
	ErrFilePathAlreadyExists = errors.New("file path already exists")
)
