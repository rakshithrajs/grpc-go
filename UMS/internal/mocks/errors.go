package mocks

import "errors"

// DbOperationError represents the result of a mocked database operation.
type DbOperationError int

const (
	// DbOpSuccess indicates the mocked database operation succeeded.
	DbOpSuccess DbOperationError = iota

	// DbOpInternalError indicates the mocked database operation failed due to an internal error.
	DbOpInternalError

	// DbOpDuplicateEmail indicates the mocked database operation failed due to a duplicate email.
	DbOpDuplicateEmail

	// DbOpDuplicatePhone indicates the mocked database operation failed due to a duplicate phone number.
	DbOpDuplicatePhone

	// DbOpDuplicateFile indicates the mocked database operation failed due to a duplicate file.
	DbOpDuplicateFile

	// DbOpNotFound indicates the mocked database operation failed because the record was not found.
	DbOpNotFound

	// DbOpRollbackFailure indicates the mocked database operation failed during rollback.
	DbOpRollbackFailure
)

// GrpcOperationError represents the result of a mocked gRPC operation.
type GrpcOperationError int

const (
	// GrpcOpSuccess indicates the mocked gRPC operation succeeded.
	GrpcOpSuccess GrpcOperationError = iota

	// GrpcOpMissingMetadata indicates the mocked gRPC operation failed due to missing metadata.
	GrpcOpMissingMetadata

	// GrpcOpMissingUserID indicates the mocked gRPC operation failed due to a missing user ID in metadata.
	GrpcOpMissingUserID

	// GrpcOpInternalError indicates the mocked gRPC operation failed due to an internal error.
	GrpcOpInternalError

	// GrpcOpFileAlreadyExists indicates the mocked gRPC operation failed because the file already exists.
	GrpcOpFileAlreadyExists

	// GrpcOpNotFound indicates the mocked gRPC operation failed because the file was not found.
	GrpcOpNotFound

	// GrpcOpRollbackFailure indicates the mocked gRPC operation failed during rollback.
	GrpcOpRollbackFailure

	// GrpcOpInvalidToken indicates the mocked gRPC operation failed because the token is invalid.
	GrpcOpInvalidToken

	// GrpcOpMissingBearerToken indicates the mocked gRPC operation failed because the bearer token is missing.
	GrpcOpMissingBearerToken
)

var (
	// ErrMissingMetadata is returned when metadata is missing from a gRPC request.
	ErrMissingMetadata = errors.New("missing metadata")

	// ErrMissingUserID is returned when the user ID is missing from gRPC metadata.
	ErrMissingUserID = errors.New("missing user id in metadata")

	// ErrFileAlreadyExists is returned when a file with the same name already exists.
	ErrFileAlreadyExists = errors.New("file already exists")

	// ErrFailedToDeleteFile is returned when a file deletion operation fails.
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// ErrFailedToGenerateToken is returned when the token generation fails.
	ErrFailedToGenerateToken = errors.New("failed to generate token")

	// ErrInvalidToken is returned when the provided token is invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrMissingBearerToken is returned when the bearer token is missing.
	ErrMissingBearerToken = errors.New("missing bearer token")
)
