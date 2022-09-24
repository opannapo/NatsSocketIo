package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"qr/dto"
)

var Qr = qr{}

type qr struct{}

func (q *qr) Create(w http.ResponseWriter, r *http.Request) (result interface{}, err error) {
	payload := dto.CreateQrRequest{}
	err = payload.Validate(r)
	if err != nil {
		log.Err(err).Send()
		return
	}

	return payload, err
}
