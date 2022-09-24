package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"qr/dto"
	"qr/service"
)

var Qr = qr{}

type qr struct{}

func (q *qr) Create(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
	reqId := getRequestID(r)
	payload := dto.CreateQrRequest{}
	err = payload.Validate(r)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	result, err = service.QrService.Create(payload)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return nil, err
	}

	return
}
