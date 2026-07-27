package mocks

type DbOperationError int

const (
	// Represents a successful database operation
	DbOpSuccess DbOperationError = iota

	// Represents a database operation that failed due to an internal error
	DbOpInternalError

	// Represents a database operation that failed due to a duplicate name
	DbOpDuplicateName

	// Represents a database operation that failed because the requested item was not found
	DbOpNotFound

	// Represents a database operation that failed due to a rollback error
	DbOpRollbackError
)
