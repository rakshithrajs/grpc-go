package mocks

import "errors"

type DbOperationError int

const (
	// Represents a successful database operation
	DbOpSuccess DbOperationError = iota

	// Represents a database operation that failed due to an internal error
	DbOpInternalError

	// Represents a database operation that failed due to a duplicate name
	DbOpDuplicateEmail

	// Represents a database operation that failed due to a duplicate phone number
	DbOpDuplicatePhone

	// Represents a database operation that failed due to a duplicate file
	DbOpDuplicateFile

	// Represents a database operation that failed because the requested record was not found
	DbOpNotFound

	// Represents a database operation that failed due to a rollback failure
	DbOpRollbackFailure
)

type GrpcOperationError int

const (
	// Represents a successful gRPC operation
	GrpcOpSuccess GrpcOperationError = iota

	// Represents a gRPC operation that failed due to missing metadata
	GrpcOpMissingMetadata

	// Represents a gRPC operation that failed due to missing user ID in metadata
	GrpcOpMissingUserID

	// Represents a gRPC operation that failed due to an internal error
	GrpcOpInternalError

	// Represents a gRPC operation that failed due to a file already existing
	GrpcOpFileAlreadyExists

	// Represents a gRPC operation that failed due to a file not being found
	GrpcOpNotFound

	// Represents a gRPC operation that failed due to a rollback failure
	GrpcOpRollbackFailure
)

var (
	// message: missing metadata
	ErrMissingMetadata = errors.New("missing metadata")

	// message: missing user id in metadata
	ErrMissingUserID = errors.New("missing user id in metadata")

	// message: file already exists
	ErrFileAlreadyExists = errors.New("file already exists")
)
