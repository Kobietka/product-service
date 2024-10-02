package v1

import "errors"

var (
	ErrorProductDoesNotExist = errors.New("PRODUCT_DOES_NOT_EXIST")
	ErrorDataNotFound        = errors.New("DATA_NOT_FOUND")
	ErrorInvalidData         = errors.New("PROVIDED_DATA_INVALID")
)

type ErrorResponse struct {
	Code string `json:"code"`
}
