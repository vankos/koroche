package urlStorage

import "fmt"

type UrlStorageErrorCode int

const (
	Unknown UrlStorageErrorCode = iota
	NotFound
	AlreadyExists
	StorageFailure
)

type UrlStorageError struct {
	Code UrlStorageErrorCode
}

func (storageError *UrlStorageError) Error() string {
	return fmt.Sprintf("UrlStorageError - %d", storageError.Code)
}
