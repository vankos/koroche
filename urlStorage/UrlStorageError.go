package urlStorage

import "fmt"

// UrlStorageErrorCode represents the type of error that occurred in the URL storage operations
type UrlStorageErrorCode int

const (
	// Unknown represents an unknown error
	Unknown UrlStorageErrorCode = iota
	// NotFound indicates that the requested URL was not found in the storage
	NotFound
	// AlreadyExists indicates that the URL already exists in the storage
	AlreadyExists
	// StorageFailure indicates a failure in the storage mechanism
	StorageFailure
)

// UrlStorageError represents an error that occurred during URL storage operations, containing an error code
type UrlStorageError struct {
	Code UrlStorageErrorCode
}

// Error returns a string representation of the UrlStorageError, including the error code
func (storageError *UrlStorageError) Error() string {
	return fmt.Sprintf("UrlStorageError - %d", storageError.Code)
}
