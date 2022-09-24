package handler

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"qr/dto"
	ierr "qr/error"
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
		return nil, err
	}

	result, err = service.QrService.Create(r.Context(), payload)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return nil, err
	}

	return
}

func (q *qr) Scan(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
	reqId := getRequestID(r)

	params := mux.Vars(r)
	paramQrID := params["id"]
	if paramQrID == "" {
		err = ierr.ErrInvalidParamQrID
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	result, err = service.QrService.Scan(r.Context(), paramQrID)
	if err != nil {
		log.Err(err).Interface("request_id", reqId).Send()
		return
	}

	return
}
