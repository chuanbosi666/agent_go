package memory

import "errors"

var (
	// ErrInvalidSessionID indicates that the provided session ID is empty or invalid.
	ErrInvalidSessionID = errors.New("session ID is required")
	// ErrDatabaseOpen indicates a failure to open the SQLite database.
	ErrDatabaseOpen = errors.New("failed to open database")
	// ErrDatabaseInit indicates a failure to initialize the database schema.
	ErrDatabaseInit = errors.New("failed to initialize database schema")
	// ErrSessionNotFound indicates that the session does not exist.
	ErrSessionNotFound = errors.New("session not found")
	// ErrInvalidItemData indicates that an item could not be serialized/deserialized.
	ErrInvalidItemData = errors.New("invalid item data")
	// ErrTransactionFailed indicates that a database transaction failed.
	ErrTransactionFailed = errors.New("database transaction failed")
	// ErrOperationFailed indicates a failure in a database operation.
	ErrOperationFailed = errors.New("database operation failed")
	// ErrDatabaseClose indicates a failure to close the database connection.
	ErrDatabaseClose = errors.New("failed to close database")
)
