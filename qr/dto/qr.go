package dto

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	ierr "qr/error"
)

type CreateQrRequest struct {
	QrId string `json:"qr_id" validate:"required"`
}

func (cqr CreateQrRequest) Validate(r *http.Request) (err error) {
	err = json.NewDecoder(r.Body).Decode(&cqr)
	if err != nil {
		log.Err(err).Send()
		return ierr.ErrInvalidRequestPayload
	}
	return
}
