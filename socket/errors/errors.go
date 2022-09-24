package errors

import (
	"socket/dto"
)

func SocketErrorAuthorization(err error) dto.SocketError {
	return dto.SocketError{
		Code:    401,
		Message: err.Error(),
	}
}

var SocketErrorInvalidHeaderXQrCodesId = dto.SocketError{Code: 400, Message: "error invalid header x-qrcodesId"}
var SocketErrorTimeToLive = dto.SocketError{Code: 403, Message: "SocketErrorTimeToLive"}
