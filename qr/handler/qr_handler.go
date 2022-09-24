package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"qr/dto"
)

func Create(w http.ResponseWriter, r *http.Request) {
	payload := dto.CreateQrRequest{}
	err := payload.Validate(r)
	if err != nil {
		log.Err(err).Send()
	}

	Response(r, w, payload, err)
}
