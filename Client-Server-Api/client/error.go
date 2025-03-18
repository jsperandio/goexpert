package client

import (
	"errors"
	"fmt"
)

var ErrExternalRequestTimeout = errors.New("external request context timeout")

type ExternalRequestError struct {
	Code int `json:"code"`
}

func NewExternalRequestError(code int) *ExternalRequestError {
	return &ExternalRequestError{
		Code: code,
	}
}

func (e *ExternalRequestError) Error() string {
	return fmt.Sprintf("external request returned error with code: %d", e.Code)
}
