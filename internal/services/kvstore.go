package services

import "fmt"

type IKeyValueStore interface {
	//set the value to storage by key
	Set(key string, value string) error
	//get the value from storage by key
	Get(key string) (string, error)
	//delete the value from storage by key
	Delete(key string) error
}

var (
	CodeKeyNotFound  = "key not found"
	CodeKeyExists    = "key already exists"
	CodeInvalidInput = "invalid input"

	KeyMaxLength   = 64
	ValueMaxLength = 1024
)

type ServiceError struct {
	Code    string
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

func validateInput(input string, maxLength int) error {
	if input == "" {
		return ServiceError{
			Code:    CodeInvalidInput,
			Message: fmt.Sprintf("input paramter %s is empty", input),
		}
	}
	if len(input) > KeyMaxLength {
		return ServiceError{
			Code:    CodeInvalidInput,
			Message: fmt.Sprintf("input paramter %s exceeded max length of %d", input, maxLength),
		}
	}
	return nil
}
