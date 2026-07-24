package mocks

type DbOperationError int

const (
	DbOpSuccess DbOperationError = iota
	DbOpInternalError
	DbOpDuplicateName
	DbOpNotFound
)
