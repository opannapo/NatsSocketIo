package error

import (
	"errors"
	"net/http"
)

var InternalErrors = NewInternalErrors()

type internalError struct {
	Code       int
	Message    string
	HttpStatus int
}

var ErrInvalidRequestPayload = errors.New("ErrInvalidRequestPayload")
var ErrInvalidCreateQrAmount = errors.New("ErrInvalidCreateQrAmoung")

func NewInternalErrors() map[error]internalError {
	return map[error]internalError{
		ErrInvalidRequestPayload: {Code: 1, Message: "Invalid Request Payload", HttpStatus: http.StatusBadRequest},
		ErrInvalidCreateQrAmount: {Code: 1, Message: "Invalid QR Amount", HttpStatus: http.StatusBadRequest},
	}
}
