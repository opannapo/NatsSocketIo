package errors

import (
	"errors"
	"socket/dto"
)

func SocketErrorAuthorization(err error) dto.SocketError {
	return dto.SocketError{
		Code:    401,
		Message: err.Error(),
	}
}

func SocketErrorInvalidHeaderXQrCodesId(err error) dto.SocketError {
	return dto.SocketError{
		Code:    400,
		Message: err.Error(),
	}
}

var ErrorDbRecordNotFound = errors.New("ErrorRecordNotFound")
var ErrorDbExec = errors.New("ErrorDbExec")
var ErrorDbQuery = errors.New("ErrorDbQuery")
var ErrorUserAccessDeniedOnlyTanky = errors.New("ErrorUserAccessDeniedOnlyTanky")
var ErrorUserAccessDeniedOnlyCompany = errors.New("ErrorUserAccessDeniedOnlyCompany")
var ErrorUserAccessDeniedOnlyDealer = errors.New("ErrorUserAccessDeniedOnlyDealer")
var ErrorUserInvalidOldPassword = errors.New("ErrorUserInvalidOldPassword")
var ErrorHashingPassword = errors.New("ErrorHashingPassword")
var ErrorActiveSessionNotFound = errors.New("ErrorActiveSessionNotFound")
var ErrorClientTankyMasterSendEmail = errors.New("ErrorClientTankyMasterSendEmail")
var ErrorEmailNotRegistered = errors.New("ErrorEmailNotRegistered")

var ErrorFieldMinimumLength = errors.New("ErrorFieldMinimumLength")
var ErrorFieldLowerCase = errors.New("ErrorFieldLowerCase")
var ErrorFieldUpperCase = errors.New("ErrorFieldUpperCase")
var ErrorFieldNumber = errors.New("ErrorFieldNumber")
var ErrorFieldSymbol = errors.New("ErrorFieldSymbol")
var ErrorFieldInteger = errors.New("ErrorFieldInteger")
var ErrorFieldRoleNameAlreadyExists = errors.New("ErrorFieldRoleNameAlreadyExists")

var ErrorInvalidPasswordVerfication = errors.New("ErrorInvalidPasswordVerfication")
