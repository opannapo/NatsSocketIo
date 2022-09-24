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
var ErrInvalidParamQrID = errors.New("ErrInvalidParamQrID")
var ErrUpdateQrScan = errors.New("ErrUpdateQrScan")
var ErrNoRowsAffected = errors.New("ErrNoRowsAffected")
var ErrQrNotFound = errors.New("ErrQrNotFound")

func NewInternalErrors() map[error]internalError {
	return map[error]internalError{
		ErrInvalidRequestPayload: {Code: 1, Message: "Invalid Request Payload", HttpStatus: http.StatusBadRequest},
		ErrInvalidCreateQrAmount: {Code: 2, Message: "Invalid QR Amount", HttpStatus: http.StatusBadRequest},
		ErrInvalidParamQrID:      {Code: 3, Message: "Invalid Param QrCode Id", HttpStatus: http.StatusBadRequest},
		ErrUpdateQrScan:          {Code: 4, Message: "Error Update QrScan", HttpStatus: http.StatusBadRequest},
		ErrNoRowsAffected:        {Code: 5, Message: "Error No RowsAffected. Maybe expired or already scanned!", HttpStatus: http.StatusNotFound},
		ErrQrNotFound:            {Code: 6, Message: "Error Not Found", HttpStatus: http.StatusNotFound},
	}
}
