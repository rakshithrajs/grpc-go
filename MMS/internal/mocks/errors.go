package mocks

// DbOperationError represents the result of a mocked database operation.
type DbOperationError int

const (
	// DbOpSuccess indicates the mocked database operation succeeded.
	DbOpSuccess DbOperationError = iota

	// DbOpInternalError indicates the mocked database operation failed due to an internal error.
	DbOpInternalError

	// DbOpDuplicateName indicates the mocked database operation failed due to a duplicate name.
	DbOpDuplicateName

	// DbOpNotFound indicates the mocked database operation failed because the record was not found.
	DbOpNotFound

	// DbOpRollbackError indicates the mocked database operation failed during rollback.
	DbOpRollbackError
)
